(function() {
    const { createApp, ref, onMounted } = Vue;

    createApp({
        setup() {
            const excludedIps = ref([]);
            const newIp = ref('');
            const loading = ref(false);
            const loaded = ref(false);
            const error = ref('');
            const success = ref('');

            function buildApiUrl() {
                const params = new URLSearchParams();
                params.set('path', '/admin/settings');
                return window.location.pathname + '?' + params.toString();
            }

            async function fetchSection(action, formData) {
                formData.set('action', action);
                const resp = await fetch(buildApiUrl(), {
                    method: 'POST',
                    body: formData
                });
                const data = await resp.json();
                if (data.status !== 'success') throw new Error(data.message || 'Request failed');
                return data.data || {};
            }

            async function loadIps() {
                loading.value = true;
                error.value = '';
                try {
                    const formData = new FormData();
                    const data = await fetchSection('list-ajax', formData);
                    excludedIps.value = data.excludedIps || [];
                } catch (e) {
                    error.value = e.message;
                } finally {
                    loading.value = false;
                    loaded.value = true;
                }
            }

            async function addIp() {
                if (!newIp.value.trim()) return;
                loading.value = true;
                error.value = '';
                success.value = '';
                try {
                    const formData = new FormData();
                    formData.set('ip_address', newIp.value.trim());
                    await fetchSection('add-ip-ajax', formData);
                    newIp.value = '';
                    await loadIps();
                    success.value = 'IP added to exclusion list';
                } catch (e) {
                    error.value = e.message;
                } finally {
                    loading.value = false;
                }
            }

            async function removeIp(ip) {
                loading.value = true;
                error.value = '';
                success.value = '';
                try {
                    const formData = new FormData();
                    formData.set('ip_address', ip);
                    await fetchSection('remove-ip-ajax', formData);
                    await loadIps();
                    success.value = 'IP removed from exclusion list';
                } catch (e) {
                    error.value = e.message;
                } finally {
                    loading.value = false;
                }
            }

            async function deleteVisitorsByIp(ip) {
                if (!confirm('Permanently delete ALL visitor records from IP ' + ip + '? This cannot be undone.')) return;
                loading.value = true;
                error.value = '';
                success.value = '';
                try {
                    const formData = new FormData();
                    formData.set('ip_address', ip);
                    const data = await fetchSection('delete-visitors-ajax', formData);
                    success.value = data.message || 'Visitor records deleted';
                    await loadIps();
                } catch (e) {
                    error.value = e.message;
                } finally {
                    loading.value = false;
                }
            }

            onMounted(() => {
                loadIps();
            });

            return {
                excludedIps, newIp, loading, loaded, error, success,
                addIp, removeIp, deleteVisitorsByIp
            };
        }
    }).mount('#settings-app');
})();
