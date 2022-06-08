const colors = new Map([
    ['White', 'white'],
    ['Yellow', 'yellow'],
    ['Orange', 'orange'],
    ['Magenta', 'magenta'],
    ['Red', 'red'],
    ['Maroon', 'maroon'],
    ['Purple', 'purple'],
    ['Dark Blue', 'darkblue'],
    ['Blue', 'blue'],
    ['Cyan', 'cyan'],
    ['Teal', 'teal'],
    ['Green', 'green'],
    ['Dark Green', 'darkgreen'],
    ['Brown', 'brown'],
]);

function ne(s) {
    return s !== undefined && s !== null && s !== "";
}

function getIcon(item, withText) {
    if (item.category === "contact" || (ne(item.team) && ne(item.role))) {
        return {uri: toUri(roleCircle(24, colors.get(item.team), '#000', item.role)), x: 12, y: 12};
    }
    if (ne(item.icon) && item.icon.startsWith("COT_MAPPING_SPOTMAP/")) {
        return {uri: toUri(circle(16, ne(item.color) ? item.color : 'green', '#000', null)), x: 8, y: 8}
    }
    if (item.type === "b") {
        return {uri: "/static/icons/b.png", x: 16, y: 16}
    }
    if (item.type === "b-m-p-w-GOTO") {
        return {uri: "/static/icons/green_flag.png", x: 6, y: 30}
    }
    if (item.type === "b-m-p-s-p-op") {
        return {uri: "/static/icons/binos.png", x: 16, y: 16}
    }
    if (item.type === "b-m-p-s-p-loc") {
        return {uri: "/static/icons/sensor_location.png", x: 16, y: 16}
    }
    if (item.type === "b-m-p-s-p-i") {
        return {uri: "/static/icons/b-m-p-s-p-i.png", x: 16, y: 16}
    }
    if (item.type === "b-m-p-a") {
        return {uri: "/static/icons/aimpoint.png", x: 16, y: 16}
    }
    if (item.category === "point") {
        return {uri: toUri(circle(16, ne(item.color) ? item.color : 'green', '#000', null)), x: 8, y: 8}
    }
    return getMilIcon(item, withText);
}

function getMilIcon(item, withText) {
    let opts = {size: 24};

    if (!ne(item.sidc)) {
        return "";
    }

    if (withText) {
        // opts['uniqueDesignation'] = item.callsign;
        if (item.speed > 0) {
            opts['speed'] = (item.speed * 3.6).toFixed(1) + " km/h";
            opts['direction'] = item.course;
        }
        if (item.sidc.charAt(2) === 'A') {
            opts['altitudeDepth'] = item.hae.toFixed(0) + " m";
        }
    }

    let symb = new ms.Symbol(item.sidc, opts);
    return {uri: symb.toDataURL(), x: symb.getAnchor().x, y: symb.getAnchor().y}
}

let app = new Vue({
    el: '#app',
    data: {
        units: new Map(),
        messages: [],
        map: null,
        ts: 0,
        locked_unit_uid: '',
        current_unit_uid: null,
        config: null,
        tools: new Map(),
        me: null,
        coords: null,
        point_num: 1,
        coord_format: "d",
        form_unit: {},
    },

    mounted() {
        this.map = L.map('map');
        let osm = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
            maxZoom: 19
        });
        let topoAttribution = 'Map data: &copy; <a href="https://openstreetmap.org/copyright">OpenStreetMap</a> contributors, <a href="http://viewfinderpanoramas.org">SRTM</a> | Map style: &copy; <a href="https://opentopomap.org">OpenTopoMap</a>, <a href="https://opentopomap.ru">OpenTopoMap.ru</a> (<a href="https://creativecommons.org/licenses/by-sa/3.0/">CC-BY-SA</a>)';
        let opentopo = L.tileLayer('https://tile-{s}.opentopomap.ru/{z}/{x}/{y}.png', {
            attribution: topoAttribution,
            maxZoom: 17
        });
        let google = L.tileLayer('http://{s}.google.com/vt/lyrs=y&x={x}&y={y}&z={z}&s=Galileo', {
            subdomains: ['mt1', 'mt2', 'mt3'],
            maxZoom: 20
        });
        let sputnik = L.tileLayer('https://{s}.tilessputnik.ru/{z}/{x}/{y}.png', {
            maxZoom: 20
        });
        osm.addTo(this.map);

        L.control.scale({metric: true}).addTo(this.map);
        L.control.layers({
            "OSM": osm,
            "Telesputnik": sputnik,
            "OpenTopoMap": opentopo,
            "Google sat": google
        }, null, {hideSingleBase: true}).addTo(this.map);

        this.renew();
        this.timer = setInterval(this.renew, 3000);
        setInterval(this.sendMyPoints, 60000);

        this.map.on('click', this.mapClick);
        this.map.on('mousemove', this.mouseMove);

        this.formFromUnit(null);
        let vm = this;
        fetch('/config')
            .then(function (response) {
                return response.json()
            })
            .then(function (data) {
                vm.config = data;
                if (vm.map != null) {
                    vm.map.setView([data.lat, data.lon], data.zoom);
                }

                if (ne(data.callsign)) {
                    vm.me = L.marker([data.lat, data.lon]);
                    vm.me.setIcon(L.icon({
                        iconUrl: "/static/icons/self.png",
                        iconAnchor: new L.Point(16, 16),
                    }));
                    vm.me.addTo(vm.map);
                }
            });
    },
    computed: {
        current_unit: function () {
            if (this.current_unit_uid != null) {
                return this.current_unit_uid && this.getCurrentUnit();
            } else {
                return null;
            }
        }
    },
    methods: {
        renew: function () {
            let vm = this;

            fetch('/units')
                .then(function (response) {
                    return response.json()
                })
                .then(function (data) {
                    let keys = new Set();

                    data.units.forEach(function (u) {
                        let oldUnit = vm.units.get(u.uid);
                        let updateMarker = false;
                        if (oldUnit === null || oldUnit === undefined) {
                            vm.units.set(u.uid, u);
                            oldUnit = u;
                            updateMarker = true;
                        } else {
                            updateMarker = (oldUnit.sidc !== u.sidc);
                            for (const k of Object.keys(u)) {
                                oldUnit[k] = u[k];
                            }
                        }
                        vm.updateMarker(oldUnit, false, updateMarker);
                        keys.add(oldUnit.uid);
                    });

                    for (const [k, v] in vm.units.entries()) {
                        if (v.parent_uid !== this.config.uid && !keys.has(k)) {
                            vm.removeUnit(k);
                        }
                    }

                    vm.messages = data.messages;
                    vm.ts += 1;
                });

            fetch('/config')
                .then(function (response) {
                    return response.json()
                })
                .then(function (data) {
                    vm.config = data;
                    if (ne(vm.me)) {
                        vm.me.setLatLng([data.lat, data.lon]);
                    }
                });

            if (this.getTool("dp1") != null) {
                let p = this.getTool("dp1").getLatLng();

                const requestOptions = {
                    method: "POST",
                    headers: {"Content-Type": "application/json"},
                    body: JSON.stringify({lat: p.lat, lon: p.lng, name: "DP1"})
                };
                fetch("/dp", requestOptions);
            }
        },
        sendMyPoints: function () {
            for (u of this.units) {
                if (u.my !== true || u.send !== true) {
                    return;
                }
                const requestOptions = {
                    method: "POST",
                    headers: {"Content-Type": "application/json"},
                    body: JSON.stringify(this.cleanUnit(u))
                };
                fetch("/add", requestOptions);
            }
        },
        updateMarker: function (item, draggable, updateIcon) {
            if (item.lon === 0 && item.lat === 0) {
                if (item.marker != null) {
                    this.map.removeLayer(item.marker);
                    item.marker = null;
                }
                return
            }

            if (ne(item.marker)) {
                if (updateIcon) {
                    let icon = getIcon(item, true);
                    item.marker.setIcon(L.icon({
                        iconUrl: icon.uri,
                        iconAnchor: new L.Point(icon.x, icon.y),
                    }));
                }
            } else {
                item.marker = L.marker([item.lat, item.lon], {draggable: draggable});
                item.marker.on('click', function (e) {
                    app.setCurrentUnitUid(item.uid, false);
                });
                if (draggable) {
                    item.marker.on('dragend', function (e) {
                        item.lat = marker.getLatLng().lat;
                        item.lon = marker.getLatLng().lng;
                    });
                }
                let icon = getIcon(item, true);
                item.marker.setIcon(L.icon({
                    iconUrl: icon.uri,
                    iconAnchor: new L.Point(icon.x, icon.y),
                }));
                item.marker.addTo(this.map);
            }

            item.marker.setLatLng([item.lat, item.lon]);
            item.marker.bindTooltip(popup(item));
            if (this.locked_unit_uid === item.uid) {
                this.map.setView([item.lat, item.lon]);
            }
        },
        removeUnit: function (uid) {
            if (!this.units.has(uid)) return;

            let item = this.units.get(uid);
            if (item.marker != null) {
                this.map.removeLayer(item.marker);
                item.marker.remove();
            }
            this.units.delete(uid);
            if (this.current_unit_uid === uid) {
                this.setCurrentUnitUid(null, false);
            }
        },
        setCurrentUnitUid: function (uid, follow) {
            if (uid != null && this.units.has(uid)) {
                this.current_unit_uid = uid;
                let u = this.units.get(uid);
                if (follow) this.mapToUnit(u);
                this.formFromUnit(u);
            } else {
                this.current_unit_uid = null;
                this.formFromUnit(null);
            }
        },
        getCurrentUnit: function () {
            if (this.current_unit_uid == null || !this.units.has(this.current_unit_uid)) return null;
            return this.units.get(this.current_unit_uid);
        },
        formFromUnit: function (u) {
            if (u == null) {
                this.form_unit = {
                    callsign: "",
                    category: "",
                    type: "",
                    subtype: "",
                    aff: "",
                    send: false,
                };
            } else {
                this.form_unit = {
                    callsign: u.callsign,
                    category: u.category,
                    type: u.type,
                    subtype: "G",
                    aff: "h",
                    send: u.send,
                };

                if (u.type.startsWith('a-')) {
                    this.form_unit.type = 'a';
                    this.form_unit.aff = u.type.substring(2, 3);
                    this.form_unit.subtype = u.type.substring(4);
                }
            }
        },
        byCategory: function (s) {
            let arr = Array.from(this.units.values()).filter(function (u) {
                return u.category === s
            });
            arr.sort(function (a, b) {
                let ua = a.callsign.toLowerCase(), ub = b.callsign.toLowerCase();
                if (ua < ub) return -1;
                if (ua > ub) return 1;
                return 0;
            });
            return this.ts && arr;
        },
        mapToUnit: function (u) {
            if (u == null) {
                return;
            }
            if (u.lat !== 0 || u.lon !== 0) {
                this.map.setView([u.lat, u.lon]);
            }
        },
        getImg: function (item) {
            return getIcon(item, false).uri;
        },
        milImg: function (item) {
            return getMilIcon(item, false).uri;
        },
        dt: function (str) {
            let d = new Date(Date.parse(str));
            return ("0" + d.getDate()).slice(-2) + "-" + ("0" + (d.getMonth() + 1)).slice(-2) + "-" +
                d.getFullYear() + " " + ("0" + d.getHours()).slice(-2) + ":" + ("0" + d.getMinutes()).slice(-2);
        },
        sp: function (v) {
            return (v * 3.6).toFixed(1);
        },
        modeIs: function (s) {
            return document.getElementById(s).checked === true;
        },
        mapClick: function (e) {
            if (this.modeIs("redx")) {
                this.addOrMove("redx", e.latlng, "/static/icons/x.png")
                return;
            }
            if (this.modeIs("dp1")) {
                this.addOrMove("dp1", e.latlng, "/static/icons/spoi_icon.png")
                return;
            }
            if (this.modeIs("point")) {
                let uid = uuidv4();
                let now = new Date();
                let stale = new Date(now);
                stale.setDate(stale.getDate() + 30);
                let u = {
                    uid: uid,
                    category: "point",
                    callsign: "point-" + this.point_num++,
                    sidc: "",
                    start_time: now,
                    last_seen: now,
                    stale_time: stale,
                    type: "b-m-p-s-m",
                    lat: e.latlng.lat,
                    lon: e.latlng.lng,
                    hae: 0,
                    speed: 0,
                    course: 0,
                    status: "",
                    text: "",
                    parent_uid: "",
                    parent_callsign: "",
                    my: true,
                }
                if (this.config != null && ne(this.config.uid)) {
                    u.parent_uid = this.config.uid;
                    u.parent_callsign = this.config.callsign;
                }
                this.units.set(uid, u);
                this.updateMarker(u, true, false);
                this.setCurrentUnitUid(u.uid, true);
            }
            if (this.modeIs("me")) {
                this.config.lat = e.latlng.lat;
                this.config.lon = e.latlng.lng;
                this.me.setLatLng(e.latlng);
                const requestOptions = {
                    method: "POST",
                    headers: {"Content-Type": "application/json"},
                    body: JSON.stringify({lat: e.latlng.lat, lon: e.latlng.lng})
                };
                fetch("/pos", requestOptions);
            }
        },
        mouseMove: function (e) {
            this.coords = e.latlng;
        },
        removeTool: function (name) {
            if (this.tools.has(name)) {
                let p = this.tools.get(name);
                this.map.removeLayer(p);
                p.remove();
                this.tools.delete(name);
                this.ts++;
            }
        },
        getTool: function (name) {
            return this.tools.get(name);
        },
        addOrMove(name, coord, icon) {
            if (this.tools.has(name)) {
                this.tools.get(name).setLatLng(coord);
            } else {
                let p = new L.marker(coord).addTo(this.map);
                if (ne(icon)) {
                    p.setIcon(L.icon({
                        iconUrl: icon,
                        iconSize: [20, 20],
                        iconAnchor: new L.Point(10, 10),
                    }));
                }
                this.tools.set(name, p);
            }
            this.ts++;
        },
        printCoordsll: function (latlng) {
            return this.printCoords(latlng.lat, latlng.lng);
        },
        printCoords: function (lat, lng) {
            return lat.toFixed(6) + "," + lng.toFixed(6);
        },
        latlng: function (lat, lon) {
            return L.latLng(lat, lon);
        },
        distBea: function (p1, p2) {
            let toRadian = Math.PI / 180;
            // haversine formula
            // bearing
            let y = Math.sin((p2.lng - p1.lng) * toRadian) * Math.cos(p2.lat * toRadian);
            let x = Math.cos(p1.lat * toRadian) * Math.sin(p2.lat * toRadian) - Math.sin(p1.lat * toRadian) * Math.cos(p2.lat * toRadian) * Math.cos((p2.lng - p1.lng) * toRadian);
            let brng = Math.atan2(y, x) * 180 / Math.PI;
            brng += brng < 0 ? 360 : 0;
            // distance
            let R = 6371000; // meters
            let deltaF = (p2.lat - p1.lat) * toRadian;
            let deltaL = (p2.lng - p1.lng) * toRadian;
            let a = Math.sin(deltaF / 2) * Math.sin(deltaF / 2) + Math.cos(p1.lat * toRadian) * Math.cos(p2.lat * toRadian) * Math.sin(deltaL / 2) * Math.sin(deltaL / 2);
            let c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
            let distance = R * c;
            return (distance < 10000 ? distance.toFixed(0) + "m " : (distance / 1000).toFixed(1) + "km ") + brng.toFixed(1) + "°T";
        },
        contactsNum: function () {
            let online = 0;
            let total = 0;
            this.units.forEach(function (u) {
                if (u.category === "contact") {
                    if (u.status === "Online") online += 1;
                    if (u.status !== "") total += 1;
                }
            })

            return online + "/" + total;
        },
        countByCategory: function (s) {
            let total = 0;
            this.units.forEach(function (u) {
                if (u.category === s) total += 1;
            })

            return total;
        },
        msgNum: function () {
            if (this.messages == null) return 0;
            return this.messages.length;
        },
        ne: function (s) {
            return s !== undefined && s !== null && s !== "";
        },
        getUnitName: function (u) {
            let res = u.callsign;
            if (u.parent_uid === this.config.uid) {
                if (u.send === true) {
                    res = "+ " + res;
                } else {
                    res = "* " + res;
                }
            }
            return res;
        },
        saveEditForm: function () {
            let u = this.getCurrentUnit();
            if (!ne(u)) return;

            u.callsign = this.form_unit.callsign;
            u.category = this.form_unit.category;

            if (this.form_unit.category === "unit") {
                u.type = ["a", this.form_unit.aff, this.form_unit.subtype].join('-');
                u.sidc = this.sidcFromType(u.type);
                u.send = this.form_unit.send;
            } else {
                u.type = this.form_unit.type;
                u.sidc = "";
                u.send = this.form_unit.send;
            }
            this.updateMarker(u, false, true);

            if (u.send === true) {
                const requestOptions = {
                    method: "POST",
                    headers: {"Content-Type": "application/json"},
                    body: JSON.stringify(this.cleanUnit(u))
                };
                fetch("/add", requestOptions);
            }
        },
        cancelEditForm: function () {
            this.formFromUnit(this.getCurrentUnit());
        },
        sidcFromType: function (s) {
            if (!s.startsWith('a-')) return "";

            let n = s.split('-');

            let sidc = 'S' + n[1];

            if (n.length > 2) {
                sidc += n[2] + 'P';
            } else {
                sidc += '-P';
            }

            if (n.length > 3) {
                for (let i = 3; i < n.length; i++) {
                    if (n[i].length > 1) {
                        break
                    }
                    sidc += n[i];
                }
            }

            if (sidc.length < 10) {
                sidc += '-'.repeat(10 - sidc.length);
            }

            return sidc.toUpperCase();
        },
        cleanUnit: function (u) {
            let res = {};

            for (const k in u) {
                if (k != 'marker' && k != 'my' && k != 'send') {
                    res[k] = u[k];
                }
            }
            return res;
        },
        deleteCurrentUnit: function () {
            if (this.current_unit_uid == null) return;
            fetch("unit/" + this.current_unit_uid, {method: "DELETE"});
            this.removeUnit(this.current_unit_uid);
        }
    },
});

function popup(item) {
    let v = '<b>' + item.callsign + '</b><br/>';
    if (ne(item.team)) v += item.team + ' ' + item.role + '<br/>';
    if (ne(item.speed) && item.speed > 0) v += 'Speed: ' + item.speed.toFixed(0) + ' m/s<br/>';
    if (item.sidc.charAt(2) === 'A') {
        v += "hae: " + item.hae.toFixed(0) + " m<br/>";
    }
    v += item.text.replaceAll('\n', '<br/>').replaceAll('; ', '<br/>');
    return v;
}

function circle(size, color, bg, text) {
    let x = Math.round(size / 2);
    let r = x - 1;

    let s = '<svg width="' + size + '" height="' + size + '" xmlns="http://www.w3.org/2000/svg"><metadata id="metadata1">image/svg+xml</metadata>';
    s += '<circle style="fill: ' + color + '; stroke: ' + bg + ';" cx="' + x + '" cy="' + x + '" r="' + r + '"/>';

    if (text != null && text !== '') {
        s += '<text x="50%" y="50%" text-anchor="middle" font-size="12px" font-family="Arial" dy=".3em">' + text + '</text>';
    }
    s += '</svg>';
    return s;
}

function roleCircle(size, color, bg, role) {
    let t = '';
    if (role === 'HQ') {
        t = 'HQ';
    } else if (role === 'Team Lead') {
        t = 'TL';
    } else if (role === 'K9') {
        t = 'K9';
    } else if (role === 'Forward Observer') {
        t = 'FO';
    } else if (role === 'Sniper') {
        t = 'S';
    } else if (role === 'Medic') {
        t = 'M';
    } else if (role === 'RTO') {
        t = 'R';
    }

    return circle(size, color, bg, t);
}

function toUri(s) {
    return encodeURI("data:image/svg+xml," + s).replaceAll("#", "%23");
}

function uuidv4() {
    return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}
