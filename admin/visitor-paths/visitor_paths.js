(function() {
    const { createApp, ref, reactive, computed, onMounted } = Vue;

    createApp({
        setup() {
            const paths = ref([]);
            const loading = ref(false);
            const loaded = ref(false);
            const error = ref('');
            const page = ref(1);
            const totalPages = ref(1);
            const pageSize = ref(10);
            const totalCount = ref(0);
            const showFilters = ref(false);

            const filters = reactive({
                range: '',
                country: '',
                device: '',
                pathContains: '',
                pathExact: '',
            });

            const exportUrl = computed(() => {
                const params = new URLSearchParams();
                params.set('path', '/admin/visitor-paths');
                params.set('action', 'export');
                if (filters.range) params.set('range', filters.range);
                if (filters.country) params.set('country', filters.country);
                if (filters.device) params.set('device', filters.device);
                if (filters.pathContains) params.set('path_contains', filters.pathContains);
                if (filters.pathExact) params.set('path_exact', filters.pathExact);
                if (page.value > 1) params.set('page', String(page.value));
                if (pageSize.value !== 10) params.set('per_page', String(pageSize.value));
                return window.location.pathname + '?' + params.toString();
            });

            const paginationSummary = computed(() => {
                if (totalCount.value === 0) return 'No paths to display';
                const from = (page.value - 1) * pageSize.value + 1;
                const to = Math.min(page.value * pageSize.value, totalCount.value);
                return 'Showing ' + from + '-' + to + ' of ' + totalCount.value + ' paths';
            });

            const pageNumbers = computed(() => {
                const pages = [];
                const max = totalPages.value;
                const cur = page.value;
                let start = Math.max(1, cur - 2);
                let end = Math.min(max, cur + 2);
                if (end - start < 4 && max > 5) {
                    if (start === 1) end = Math.min(max, 5);
                    else if (end === max) start = Math.max(1, max - 4);
                }
                for (let i = start; i <= end; i++) pages.push(i);
                return pages;
            });

            function buildApiUrl() {
                const params = new URLSearchParams();
                params.set('path', '/admin/visitor-paths');
                return window.location.pathname + '?' + params.toString();
            }

            async function fetchList() {
                loading.value = true;
                error.value = '';
                try {
                    const formData = new FormData();
                    formData.set('action', 'list-ajax');
                    formData.set('page', String(page.value));
                    formData.set('per_page', String(pageSize.value));
                    if (filters.range) formData.set('range', filters.range);
                    if (filters.country) formData.set('country', filters.country);
                    if (filters.device) formData.set('device', filters.device);
                    if (filters.pathContains) formData.set('path_contains', filters.pathContains);
                    if (filters.pathExact) formData.set('path_exact', filters.pathExact);

                    const resp = await fetch(buildApiUrl(), { method: 'POST', body: formData });
                    const data = await resp.json();
                    if (data.status !== 'success') throw new Error(data.message || 'Request failed');

                    const d = data.data || {};
                    paths.value = d.paths || [];
                    totalPages.value = d.totalPages || 1;
                    totalCount.value = d.totalCount || 0;
                } catch (e) {
                    error.value = e.message;
                } finally {
                    loading.value = false;
                    loaded.value = true;
                }
            }

            function applyFilters() {
                page.value = 1;
                showFilters.value = false;
                fetchList();
            }

            function clearFilters() {
                filters.range = '';
                filters.country = '';
                filters.device = '';
                filters.pathContains = '';
                filters.pathExact = '';
                page.value = 1;
                fetchList();
            }

            function changePage(p) {
                if (p < 1 || p > totalPages.value || p === page.value) return;
                page.value = p;
                fetchList();
            }

            function changePageSize(size) {
                if (size === pageSize.value) return;
                pageSize.value = size;
                page.value = 1;
                fetchList();
            }

            function rangeLabel(value) {
                const labels = { '24h': 'Last 24 Hours', 'today': 'Today', '7d': 'Last 7 Days', '30d': 'Last 30 Days' };
                return labels[value] || value;
            }

            function capitalize(s) {
                return s.charAt(0).toUpperCase() + s.slice(1);
            }

            onMounted(() => {
                const urlParams = new URLSearchParams(window.location.search);
                const range = urlParams.get('range');
                const country = urlParams.get('country');
                const device = urlParams.get('device');
                const pc = urlParams.get('path_contains');
                const pe = urlParams.get('path_exact');
                const p = urlParams.get('page');
                const pp = urlParams.get('per_page');
                if (range) filters.range = range;
                if (country) filters.country = country;
                if (device) filters.device = device;
                if (pc) filters.pathContains = pc;
                if (pe) filters.pathExact = pe;
                if (p) page.value = parseInt(p) || 1;
                if (pp) pageSize.value = parseInt(pp) || 10;
                fetchList();
            });

            return {
                paths, loading, loaded, error,
                page, totalPages, pageSize, totalCount,
                showFilters, filters,
                exportUrl, paginationSummary, pageNumbers,
                fetchList, applyFilters, clearFilters, changePage, changePageSize,
                rangeLabel, capitalize,
            };
        }
    }).mount('#visitor-paths-app');
})();
