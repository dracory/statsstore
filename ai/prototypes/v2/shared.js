/* StatsStore Prototype v2 — Shared data, helpers, toast system */

const SS = {
    navItems: [
        { title: 'Dashboard', href: 'index.html', icon: 'bi-speedometer2' },
        { title: 'Visitor Activity', href: 'visitor-activity.html', icon: 'bi-people' },
        { title: 'Visitor Paths', href: 'visitor-paths.html', icon: 'bi-signpost-split' },
        { title: 'Page Views', href: 'page-view-activity.html', icon: 'bi-eye' },
        { title: 'Settings', href: 'settings.html', icon: 'bi-gear' },
    ],

    periodOptions: [
        { value: 'today', label: 'Today' },
        { value: 'yesterday', label: 'Yesterday' },
        { value: 'last-7-days', label: 'Last 7 Days' },
        { value: 'this-week', label: 'This Week' },
        { value: 'last-week', label: 'Last Week' },
        { value: 'this-month', label: 'This Month' },
        { value: 'last-month', label: 'Last Month' },
    ],

    statCards: [
        { title: 'Page Views', value: '12,847', icon: 'bi-eye', color: 'primary' },
        { title: 'Unique Visitors', value: '8,234', icon: 'bi-person', color: 'success' },
        { title: 'First Visits', value: '5,621', icon: 'bi-person-plus', color: 'info' },
        { title: 'Returning', value: '2,613', icon: 'bi-person-check', color: 'warning' },
        { title: 'Bounce Rate', value: '42.3%', icon: 'bi-arrow-return-left', color: 'danger' },
        { title: 'Avg. Duration', value: '3m 42s', icon: 'bi-clock', color: 'secondary' },
    ],

    comparisonRows: [
        { label: 'Page Views', current: '12,847', previous: '11,203', change: 14.7, inverted: false },
        { label: 'Unique Visitors', current: '8,234', previous: '7,891', change: 4.3, inverted: false },
        { label: 'First Visits', current: '5,621', previous: '5,102', change: 10.2, inverted: false },
        { label: 'Returning Visits', current: '2,613', previous: '2,789', change: -6.3, inverted: false },
        { label: 'Bounce Rate', current: '42.3%', previous: '38.1%', change: 11.0, inverted: true },
        { label: 'Avg. Duration', current: '3m 42s', previous: '4m 12s', change: -11.9, inverted: true },
    ],

    dailyStats: [
        { date: 'Mon, Jul 21', totalVisits: 1842, uniqueVisits: 1203, firstVisits: 812, returnVisits: 391 },
        { date: 'Tue, Jul 22', totalVisits: 2103, uniqueVisits: 1456, firstVisits: 987, returnVisits: 469 },
        { date: 'Wed, Jul 23', totalVisits: 1956, uniqueVisits: 1289, firstVisits: 845, returnVisits: 444 },
        { date: 'Thu, Jul 24', totalVisits: 2241, uniqueVisits: 1567, firstVisits: 1023, returnVisits: 544 },
        { date: 'Fri, Jul 25', totalVisits: 2578, uniqueVisits: 1789, firstVisits: 1156, returnVisits: 633 },
        { date: 'Sat, Jul 26', totalVisits: 1287, uniqueVisits: 876, firstVisits: 534, returnVisits: 342 },
        { date: 'Sun, Jul 27', totalVisits: 840, uniqueVisits: 521, firstVisits: 264, returnVisits: 257 },
    ],

    chartLabels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    chartUniqueVisits: [1203, 1456, 1289, 1567, 1789, 876, 521],
    chartTotalVisits: [1842, 2103, 1956, 2241, 2578, 1287, 840],

    trafficCards: [
        {
            title: 'Countries',
            valueLabel: 'Visits',
            tabs: [
                { label: 'Top 10', entries: [
                    { label: 'United States', sessions: 3842 },
                    { label: 'Germany', sessions: 2103 },
                    { label: 'United Kingdom', sessions: 1567 },
                    { label: 'France', sessions: 1289 },
                    { label: 'Canada', sessions: 987 },
                    { label: 'Netherlands', sessions: 845 },
                    { label: 'Australia', sessions: 678 },
                    { label: 'Spain', sessions: 534 },
                    { label: 'Italy', sessions: 412 },
                    { label: 'Sweden', sessions: 289 },
                ]},
                { label: 'All', entries: [
                    { label: 'United States', sessions: 3842 },
                    { label: 'Germany', sessions: 2103 },
                    { label: 'United Kingdom', sessions: 1567 },
                    { label: 'France', sessions: 1289 },
                    { label: 'Canada', sessions: 987 },
                    { label: 'Netherlands', sessions: 845 },
                    { label: 'Australia', sessions: 678 },
                    { label: 'Spain', sessions: 534 },
                    { label: 'Italy', sessions: 412 },
                    { label: 'Sweden', sessions: 289 },
                    { label: 'Other', sessions: 564 },
                ]},
            ],
        },
        {
            title: 'Browsers',
            valueLabel: 'Sessions',
            tabs: [
                { label: 'Top 5', entries: [
                    { label: 'Chrome', sessions: 6842 },
                    { label: 'Safari', sessions: 2103 },
                    { label: 'Firefox', sessions: 1289 },
                    { label: 'Edge', sessions: 987 },
                    { label: 'Opera', sessions: 345 },
                ]},
                { label: 'All', entries: [
                    { label: 'Chrome', sessions: 6842 },
                    { label: 'Safari', sessions: 2103 },
                    { label: 'Firefox', sessions: 1289 },
                    { label: 'Edge', sessions: 987 },
                    { label: 'Opera', sessions: 345 },
                    { label: 'Other', sessions: 281 },
                ]},
            ],
        },
        {
            title: 'Operating Systems',
            valueLabel: 'Sessions',
            tabs: [
                { label: 'Top 5', entries: [
                    { label: 'Windows', sessions: 5234 },
                    { label: 'macOS', sessions: 2891 },
                    { label: 'Android', sessions: 2156 },
                    { label: 'iOS', sessions: 1789 },
                    { label: 'Linux', sessions: 567 },
                ]},
                { label: 'All', entries: [
                    { label: 'Windows', sessions: 5234 },
                    { label: 'macOS', sessions: 2891 },
                    { label: 'Android', sessions: 2156 },
                    { label: 'iOS', sessions: 1789 },
                    { label: 'Linux', sessions: 567 },
                    { label: 'Other', sessions: 210 },
                ]},
            ],
        },
        {
            title: 'Devices',
            valueLabel: 'Sessions',
            tabs: [
                { label: 'All', entries: [
                    { label: 'Desktop', sessions: 7456 },
                    { label: 'Mobile', sessions: 3945 },
                    { label: 'Tablet', sessions: 892 },
                    { label: 'Bot', sessions: 554 },
                ]},
            ],
        },
    ],

    heatmap: {
        days: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
        slots: ['00:00', '03:00', '06:00', '09:00', '12:00', '15:00', '18:00', '21:00'],
        intensities: [
            [0, 0, 1, 2, 1, 1, 0, 0],
            [0, 0, 1, 3, 2, 2, 1, 0],
            [0, 0, 1, 3, 3, 2, 1, 0],
            [0, 0, 1, 4, 3, 3, 2, 1],
            [0, 1, 2, 4, 5, 4, 3, 1],
            [1, 1, 1, 2, 3, 2, 2, 1],
            [0, 0, 0, 1, 2, 1, 1, 0],
            [0, 0, 0, 1, 1, 0, 0, 0],
        ],
    },

    visitors: [
        { id: 'v-001', ip: '192.168.1.1', country: 'United States', countryCode: 'us', path: '/home', device: 'Desktop', browser: 'Chrome', os: 'Windows', createdAt: '2025-07-26 14:23:01', sessions: 3 },
        { id: 'v-002', ip: '10.0.0.2', country: 'Germany', countryCode: 'de', path: '/products/1', device: 'Mobile', browser: 'Safari', os: 'iOS', createdAt: '2025-07-26 14:21:45', sessions: 1 },
        { id: 'v-003', ip: '172.16.0.3', country: 'United Kingdom', countryCode: 'gb', path: '/about', device: 'Desktop', browser: 'Firefox', os: 'Linux', createdAt: '2025-07-26 14:18:22', sessions: 5 },
        { id: 'v-004', ip: '192.168.1.4', country: 'France', countryCode: 'fr', path: '/blog/post-1', device: 'Tablet', browser: 'Chrome', os: 'Android', createdAt: '2025-07-26 14:15:10', sessions: 2 },
        { id: 'v-005', ip: '10.0.0.5', country: 'Canada', countryCode: 'ca', path: '/contact', device: 'Desktop', browser: 'Edge', os: 'Windows', createdAt: '2025-07-26 14:12:33', sessions: 1 },
        { id: 'v-006', ip: '172.16.0.6', country: 'Netherlands', countryCode: 'nl', path: '/home', device: 'Mobile', browser: 'Chrome', os: 'Android', createdAt: '2025-07-26 14:08:55', sessions: 4 },
        { id: 'v-007', ip: '192.168.1.7', country: 'Australia', countryCode: 'au', path: '/products/2', device: 'Desktop', browser: 'Safari', os: 'macOS', createdAt: '2025-07-26 14:05:12', sessions: 2 },
        { id: 'v-008', ip: '10.0.0.8', country: 'Spain', countryCode: 'es', path: '/pricing', device: 'Mobile', browser: 'Firefox', os: 'iOS', createdAt: '2025-07-26 14:01:48', sessions: 1 },
        { id: 'v-009', ip: '172.16.0.9', country: 'Italy', countryCode: 'it', path: '/blog/post-2', device: 'Desktop', browser: 'Chrome', os: 'Windows', createdAt: '2025-07-26 13:58:30', sessions: 3 },
        { id: 'v-010', ip: '192.168.1.10', country: 'Sweden', countryCode: 'se', path: '/home', device: 'Mobile', browser: 'Chrome', os: 'iOS', createdAt: '2025-07-26 13:55:15', sessions: 1 },
        { id: 'v-011', ip: '10.0.0.11', country: 'United States', countryCode: 'us', path: '/products/3', device: 'Desktop', browser: 'Edge', os: 'Windows', createdAt: '2025-07-26 13:52:01', sessions: 6 },
        { id: 'v-012', ip: '172.16.0.12', country: 'Germany', countryCode: 'de', path: '/about', device: 'Mobile', browser: 'Safari', os: 'Android', createdAt: '2025-07-26 13:48:44', sessions: 2 },
    ],

    visitorPaths: [
        { id: 'p-001', path: '/home', visits: 3842, uniqueVisitors: 2103, avgDuration: '2m 15s', device: 'Desktop', sessions: 1 },
        { id: 'p-002', path: '/products/1', visits: 2156, uniqueVisitors: 1289, avgDuration: '3m 42s', device: 'Mobile', sessions: 2 },
        { id: 'p-003', path: '/blog/post-1', visits: 1789, uniqueVisitors: 987, avgDuration: '5m 12s', device: 'Desktop', sessions: 1 },
        { id: 'p-004', path: '/about', visits: 1456, uniqueVisitors: 845, avgDuration: '1m 58s', device: 'Desktop', sessions: 1 },
        { id: 'p-005', path: '/pricing', visits: 1289, uniqueVisitors: 712, avgDuration: '2m 33s', device: 'Mobile', sessions: 1 },
        { id: 'p-006', path: '/contact', visits: 987, uniqueVisitors: 534, avgDuration: '1m 12s', device: 'Desktop', sessions: 1 },
        { id: 'p-007', path: '/blog/post-2', visits: 845, uniqueVisitors: 412, avgDuration: '4m 05s', device: 'Tablet', sessions: 1 },
        { id: 'p-008', path: '/products/2', visits: 678, uniqueVisitors: 389, avgDuration: '3m 18s', device: 'Mobile', sessions: 2 },
        { id: 'p-009', path: '/blog/post-3', visits: 534, uniqueVisitors: 289, avgDuration: '6m 22s', device: 'Desktop', sessions: 1 },
        { id: 'p-010', path: '/products/3', visits: 412, uniqueVisitors: 234, avgDuration: '2m 48s', device: 'Desktop', sessions: 3 },
    ],

    pageViews: [
        { id: 'pv-001', path: '/home', visitorId: 'v-001', ip: '192.168.1.1', country: 'United States', countryCode: 'us', device: 'Desktop', browser: 'Chrome', os: 'Windows', timestamp: '2025-07-26 14:23:01' },
        { id: 'pv-002', path: '/products/1', visitorId: 'v-002', ip: '10.0.0.2', country: 'Germany', countryCode: 'de', device: 'Mobile', browser: 'Safari', os: 'iOS', timestamp: '2025-07-26 14:21:45' },
        { id: 'pv-003', path: '/about', visitorId: 'v-003', ip: '172.16.0.3', country: 'United Kingdom', countryCode: 'gb', device: 'Desktop', browser: 'Firefox', os: 'Linux', timestamp: '2025-07-26 14:18:22' },
        { id: 'pv-004', path: '/blog/post-1', visitorId: 'v-004', ip: '192.168.1.4', country: 'France', countryCode: 'fr', device: 'Tablet', browser: 'Chrome', os: 'Android', timestamp: '2025-07-26 14:15:10' },
        { id: 'pv-005', path: '/contact', visitorId: 'v-005', ip: '10.0.0.5', country: 'Canada', countryCode: 'ca', device: 'Desktop', browser: 'Edge', os: 'Windows', timestamp: '2025-07-26 14:12:33' },
        { id: 'pv-006', path: '/home', visitorId: 'v-006', ip: '172.16.0.6', country: 'Netherlands', countryCode: 'nl', device: 'Mobile', browser: 'Chrome', os: 'Android', timestamp: '2025-07-26 14:08:55' },
        { id: 'pv-007', path: '/products/2', visitorId: 'v-007', ip: '192.168.1.7', country: 'Australia', countryCode: 'au', device: 'Desktop', browser: 'Safari', os: 'macOS', timestamp: '2025-07-26 14:05:12' },
        { id: 'pv-008', path: '/pricing', visitorId: 'v-008', ip: '10.0.0.8', country: 'Spain', countryCode: 'es', device: 'Mobile', browser: 'Firefox', os: 'iOS', timestamp: '2025-07-26 14:01:48' },
        { id: 'pv-009', path: '/blog/post-2', visitorId: 'v-009', ip: '172.16.0.9', country: 'Italy', countryCode: 'it', device: 'Desktop', browser: 'Chrome', os: 'Windows', timestamp: '2025-07-26 13:58:30' },
        { id: 'pv-010', path: '/home', visitorId: 'v-010', ip: '192.168.1.10', country: 'Sweden', countryCode: 'se', device: 'Mobile', browser: 'Chrome', os: 'iOS', timestamp: '2025-07-26 13:55:15' },
    ],

    blockedIps: [
        { ip: '185.220.101.1', addedAt: '2025-07-25 10:23:00' },
        { ip: '194.165.16.1', addedAt: '2025-07-24 18:45:00' },
        { ip: '212.193.30.1', addedAt: '2025-07-23 09:12:00' },
    ],

    /* ===== Helpers ===== */

    heatmapColor(level) {
        const colors = ['#1c2333', '#12e198', '#14c987', '#17b176', '#1a9a65', '#1f8254'];
        return colors[level] || colors[0];
    },

    deviceBadgeClass(device) {
        const map = { Desktop: 'bg-primary', Mobile: 'bg-success', Tablet: 'bg-info', Bot: 'bg-danger' };
        return map[device] || 'bg-secondary';
    },

    totals() {
        return this.dailyStats.reduce((acc, r) => {
            acc.totalVisits += r.totalVisits;
            acc.uniqueVisits += r.uniqueVisits;
            acc.firstVisits += r.firstVisits;
            acc.returnVisits += r.returnVisits;
            return acc;
        }, { totalVisits: 0, uniqueVisits: 0, firstVisits: 0, returnVisits: 0 });
    },

    isValidIp(ip) {
        const v4 = /^(\d{1,3}\.){3}\d{1,3}$/;
        if (v4.test(ip)) return ip.split('.').every(o => parseInt(o) >= 0 && parseInt(o) <= 255);
        return /^[0-9a-f:]+$/i.test(ip);
    },

    /* ===== Toast System ===== */

    toasts: [],

    toast(message, type = 'info', duration = 3500) {
        const id = Date.now() + Math.random();
        this.toasts.push({ id, message, type });
        setTimeout(() => {
            const i = this.toasts.findIndex(t => t.id === id);
            if (i >= 0) this.toasts.splice(i, 1);
            this._renderToasts();
        }, duration);
        this._renderToasts();
    },

    _renderToasts() {
        let container = document.getElementById('ss-toast-container');
        if (!container) {
            container = document.createElement('div');
            container.id = 'ss-toast-container';
            container.className = 'ss-toast-container';
            document.body.appendChild(container);
        }
        const icons = { success: 'bi-check-circle-fill', error: 'bi-x-circle-fill', warning: 'bi-exclamation-triangle-fill', info: 'bi-info-circle-fill' };
        container.innerHTML = this.toasts.map(t =>
            `<div class="ss-toast ${t.type}"><i class="bi ${icons[t.type] || icons.info}"></i><span>${t.message}</span></div>`
        ).join('');
    },

    /* ===== Theme ===== */

    initTheme() {
        const saved = localStorage.getItem('ss-theme');
        if (saved) document.documentElement.setAttribute('data-theme', saved);
    },

    toggleTheme() {
        const current = document.documentElement.getAttribute('data-theme');
        const next = current === 'dark' ? 'light' : 'dark';
        if (next === 'light') document.documentElement.removeAttribute('data-theme');
        else document.documentElement.setAttribute('data-theme', 'dark');
        localStorage.setItem('ss-theme', next);
    },

    getThemeIcon() {
        return document.documentElement.getAttribute('data-theme') === 'dark' ? 'bi-sun' : 'bi-moon-stars';
    },

    /* ===== Shell Setup ===== */

    initShell() {
        this.initTheme();
    },
};

SS.initTheme();
