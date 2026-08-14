const { createApp, ref, onMounted } = Vue;

createApp({
    setup() {
        const details = ref({});
        const paths = ref([]);
        const visitCount = ref(0);
        const isBot = ref(false);
        const isThreat = ref(false);
        const loading = ref(true);
        const flagging = ref(false);
        const flaggingThreat = ref(false);
        const removing = ref(false);

        const loadDetails = async () => {
            try {
                const response = await fetch(urlLoadPaths, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ ip: '__IP__' }),
                });
                const data = await response.json();
                if (data.status === 'success') {
                    details.value = data.data.details || {};
                    paths.value = data.data.paths || [];
                    visitCount.value = data.data.visit_count || 0;
                    isBot.value = data.data.is_bot || false;
                    isThreat.value = data.data.is_threat || false;
                } else {
                    Notiflix.Notify.failure(data.message || 'Failed to load IP details');
                }
            } catch (e) {
                Notiflix.Notify.failure('Failed to load IP details');
            } finally {
                loading.value = false;
            }
        };

        const confirmFlagBot = () => {
            Notiflix.Confirm.show(
                'Flag as Bot',
                'Are you sure you want to flag ' + (details.value.ip || '__IP__') + ' as a bot? This IP will be excluded from visitor statistics.',
                'Yes, flag it',
                'Cancel',
                () => flagBot(),
                () => {}
            );
        };

        const flagBot = async () => {
            flagging.value = true;
            try {
                const response = await fetch(urlFlagBot, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ ip: '__IP__' }),
                });
                const data = await response.json();
                if (data.status === 'success') {
                    isBot.value = true;
                    Notiflix.Notify.success('IP flagged as bot');
                } else {
                    Notiflix.Notify.failure(data.message || 'Failed to flag IP');
                }
            } catch (e) {
                Notiflix.Notify.failure('Failed to flag IP');
            } finally {
                flagging.value = false;
            }
        };

        const confirmFlagThreat = () => {
            Notiflix.Confirm.show(
                'Flag as Threat',
                'Are you sure you want to flag ' + (details.value.ip || '__IP__') + ' as a threat? All visitor records from this IP will be marked as threat=yes.',
                'Yes, flag it',
                'Cancel',
                () => flagThreat(),
                () => {}
            );
        };

        const flagThreat = async () => {
            flaggingThreat.value = true;
            try {
                const response = await fetch(urlFlagThreat, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ ip: '__IP__' }),
                });
                const data = await response.json();
                if (data.status === 'success') {
                    isThreat.value = true;
                    Notiflix.Notify.success('IP flagged as threat');
                } else {
                    Notiflix.Notify.failure(data.message || 'Failed to flag IP');
                }
            } catch (e) {
                Notiflix.Notify.failure('Failed to flag IP');
            } finally {
                flaggingThreat.value = false;
            }
        };

        onMounted(() => {
            loadDetails();
        });

        const confirmRemoveEntries = () => {
            Notiflix.Confirm.show(
                'Remove Entries',
                'Are you sure you want to permanently delete ALL visitor records and path history for ' + (details.value.ip || '__IP__') + '? This action cannot be undone.',
                'Yes, delete all',
                'Cancel',
                () => removeEntries(),
                () => {}
            );
        };

        const removeEntries = async () => {
            removing.value = true;
            try {
                const response = await fetch(urlRemoveEntries, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ ip: '__IP__' }),
                });
                const data = await response.json();
                if (data.status === 'success') {
                    Notiflix.Notify.success(data.message || 'Entries removed');
                    // Clear the local view since the records no longer exist.
                    details.value = {};
                    paths.value = [];
                    visitCount.value = 0;
                    isBot.value = false;
                    isThreat.value = false;
                } else {
                    Notiflix.Notify.failure(data.message || 'Failed to remove entries');
                }
            } catch (e) {
                Notiflix.Notify.failure('Failed to remove entries');
            } finally {
                removing.value = false;
            }
        };

        return {
            details,
            paths,
            visitCount,
            isBot,
            isThreat,
            loading,
            flagging,
            flaggingThreat,
            removing,
            confirmFlagBot,
            confirmFlagThreat,
            confirmRemoveEntries,
        };
    }
}).mount('#ip-details-app');
