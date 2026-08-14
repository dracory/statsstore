const { createApp, ref, onMounted } = Vue;

const conditionOptions = [
    { value: 'ip',          label: 'IP Address',        operators: ['equals'],    inputType: 'text',   placeholder: 'e.g. 192.168.1.1' },
    { value: 'country',     label: 'Country',           operators: ['equals'],    inputType: 'text',   placeholder: 'ISO2 code (e.g. GB)' },
    { value: 'device_type', label: 'Device Type',       operators: ['equals'],    inputType: 'select', placeholder: '' },
    { value: 'path_contains', label: 'Path contains',   operators: ['contains'],  inputType: 'text',   placeholder: 'e.g. /admin' },
    { value: 'path_exact',  label: 'Path equals',       operators: ['equals'],    inputType: 'text',   placeholder: 'e.g. [GET] /' },
    { value: 'browser',     label: 'Browser',           operators: ['equals'],    inputType: 'text',   placeholder: 'e.g. Chrome' },
    { value: 'os',          label: 'Operating System',  operators: ['equals'],    inputType: 'text',   placeholder: 'e.g. Windows' },
    { value: 'date_from',   label: 'Created after',     operators: ['equals'],    inputType: 'date',   placeholder: '' },
    { value: 'date_to',     label: 'Created before',    operators: ['equals'],    inputType: 'date',   placeholder: '' },
    { value: 'is_bot',      label: 'Is Bot',            operators: ['equals'],    inputType: 'select', placeholder: '' },
];

const deviceTypes = ['desktop', 'mobile', 'tablet', 'bot'];

const opLabels = { equals: '=', contains: 'contains' };

function getFieldOpt(value) {
    return conditionOptions.find(o => o.value === value);
}

createApp({
    setup() {
        const sessions = ref([]);
        const total = ref(0);
        const page = ref(1);
        const perPage = ref(25);
        const totalPages = ref(0);

        // Applied conditions (sent to the server)
        const conditions = ref([]);

        // Modal state
        const showFilterModal = ref(false);
        const modalConditions = ref([]);

        const openFilterModal = () => {
            // Clone current conditions into the modal editor
            modalConditions.value = conditions.value.map(c => ({ ...c }));
            if (modalConditions.value.length === 0) {
                modalConditions.value.push({ field: '', operator: 'equals', value: '' });
            }
            showFilterModal.value = true;
        };

        const addModalCondition = () => {
            modalConditions.value.push({ field: '', operator: 'equals', value: '' });
        };

        const removeModalCondition = (i) => {
            modalConditions.value.splice(i, 1);
        };

        const clearModalConditions = () => {
            modalConditions.value = [];
        };

        const onFieldChange = (c) => {
            const opt = getFieldOpt(c.field);
            if (opt) {
                c.operator = opt.operators[0];
            }
            c.value = '';
        };

        const applyConditions = () => {
            // Filter out empty conditions
            conditions.value = modalConditions.value.filter(c => c.field && c.value.trim());
            showFilterModal.value = false;
            page.value = 1;
            loadSessions();
            updateURL();
        };

        const clearConditions = () => {
            conditions.value = [];
            page.value = 1;
            loadSessions();
            updateURL();
        };

        const removeCondition = (i) => {
            conditions.value.splice(i, 1);
            page.value = 1;
            loadSessions();
            updateURL();
        };

        const loadSessions = async () => {
            try {
                const response = await fetch(urlLoadSessions, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        page: page.value,
                        per_page: perPage.value,
                        conditions: conditions.value,
                    }),
                });
                const data = await response.json();
                if (data.status === 'success') {
                    sessions.value = data.data.sessions || [];
                    total.value = data.data.total || 0;
                    totalPages.value = data.data.total_pages || 0;
                } else {
                    Notiflix.Notify.failure(data.message || 'Failed to load sessions');
                }
            } catch (error) {
                Notiflix.Notify.failure('Failed to load sessions');
            }
        };

        const prevPage = () => {
            if (page.value > 1) {
                page.value--;
                loadSessions();
                updateURL();
            }
        };

        const nextPage = () => {
            if (page.value < totalPages.value) {
                page.value++;
                loadSessions();
                updateURL();
            }
        };

        // Encode/decode filter conditions to/from URL query string so that
        // saved links restore the active filters.
        const updateURL = () => {
            const params = new URLSearchParams();
            if (conditions.value.length > 0) {
                params.set('filters', JSON.stringify(conditions.value));
            }
            if (page.value > 1) {
                params.set('page', String(page.value));
            }
            const qs = params.toString();
            const newURL = qs ? window.location.pathname + '?' + qs : window.location.pathname;
            history.replaceState(null, '', newURL);
        };

        const loadFromURL = () => {
            const params = new URLSearchParams(window.location.search);
            const filtersRaw = params.get('filters');
            if (filtersRaw) {
                try {
                    const parsed = JSON.parse(filtersRaw);
                    if (Array.isArray(parsed)) {
                        conditions.value = parsed.filter(c => c && c.field && c.value);
                    }
                } catch (e) { /* ignore malformed */ }
            }
            const pageRaw = params.get('page');
            if (pageRaw) {
                const p = parseInt(pageRaw, 10);
                if (p > 0) page.value = p;
            }
        };

        // Helper functions for the template
        const fieldLabel = (v) => { const o = getFieldOpt(v); return o ? o.label : v; };
        const opLabel = (op) => opLabels[op] || op;
        const getFieldOperators = (v) => { const o = getFieldOpt(v); return o ? o.operators : ['equals']; };
        const getFieldInputType = (v) => { const o = getFieldOpt(v); return o ? o.inputType : 'text'; };
        const getFieldPlaceholder = (v) => { const o = getFieldOpt(v); return o ? o.placeholder : ''; };
        const ipDetailsUrl = (ip) => urlIPDetailsBase + '&ip=' + encodeURIComponent(ip);
        // pathUrl strips the "[METHOD] " prefix from a stored path (e.g. "[GET] /foo")
        // and builds an absolute URL so the link can open in a new tab.
        const pathUrl = (path) => {
            if (!path) return '#';
            const stripped = String(path).replace(/^\[[^\]]+\]\s*/, '');
            if (!stripped.startsWith('/')) return stripped;
            return window.location.origin + stripped;
        };

        onMounted(() => {
            loadFromURL();
            loadSessions();
        });

        return {
            sessions, total, page, perPage, totalPages,
            conditions, showFilterModal, modalConditions,
            conditionOptions, deviceTypes,
            openFilterModal, addModalCondition, removeModalCondition,
            clearModalConditions, onFieldChange, applyConditions,
            clearConditions, removeCondition,
            loadSessions, prevPage, nextPage,
            fieldLabel, opLabel, getFieldOperators, getFieldInputType, getFieldPlaceholder,
            ipDetailsUrl,
            pathUrl,
        };
    }
}).mount('#stats-sessions-app');
