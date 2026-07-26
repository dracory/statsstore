(function() {
    const { createApp, ref, reactive, computed, onMounted, onUnmounted, watch, nextTick } = Vue;

    createApp({
        setup() {
            const selectedPeriod = ref('this-week');
            const periodOptions = ref([]);
            const liveVisitorCount = ref(0);
            const statCards = ref([]);
            const comparisonRows = ref([]);
            const previousPeriodLabel = ref('');
            const dailyStats = ref([]);
            const totals = ref({});
            const trafficCards = ref([]);
            const heatmap = ref({ days: [], slots: [], intensities: [] });
            const metrics = ref(['Unique Visitors', 'Pageviews', 'Sessions', 'Bounce Rate', 'Pages per Session', 'Session Duration']);
            const selectedMetric = ref('Unique Visitors');
            const chartCanvas = ref(null);
            let chartInstance = null;
            let chartType = ref('bar');

            // Per-section loading states
            const loadingOverview = ref(false);
            const loadingComparison = ref(false);
            const loadingDaily = ref(false);
            const loadingTraffic = ref(false);
            const loadingHeatmap = ref(false);

            // Per-section error states
            const overviewError = ref('');
            const comparisonError = ref('');
            const dailyError = ref('');
            const trafficError = ref('');
            const heatmapError = ref('');

            const exportUrl = computed(() => {
                return window.location.pathname + '?path=/admin/home&action=export&period=' + selectedPeriod.value;
            });

            const visitorActivityUrl = computed(() => {
                return window.location.pathname + '?path=/admin/visitor-activity';
            });

            const visitorPathsUrl = computed(() => {
                return window.location.pathname + '?path=/admin/visitor-paths';
            });

            function buildApiUrl() {
                const params = new URLSearchParams();
                params.set('path', '/admin/home');
                return window.location.pathname + '?' + params.toString();
            }

            async function fetchSection(action, period) {
                const formData = new FormData();
                formData.set('action', action);
                if (period) formData.set('period', period);
                const resp = await fetch(buildApiUrl(), {
                    method: 'POST',
                    body: formData
                });
                const data = await resp.json();
                if (data.status !== 'success') throw new Error(data.message || 'Request failed');
                return data.data || {};
            }

            async function fetchOverview(period) {
                loadingOverview.value = true;
                overviewError.value = '';
                try {
                    const data = await fetchSection('overview-ajax', period);
                    periodOptions.value = data.periodOptions || [];
                    selectedPeriod.value = data.selectedPeriod || 'this-week';
                    liveVisitorCount.value = data.liveVisitorCount || 0;
                } catch (e) {
                    overviewError.value = e.message;
                } finally {
                    loadingOverview.value = false;
                }
            }

            async function fetchComparison(period) {
                loadingComparison.value = true;
                comparisonError.value = '';
                try {
                    const data = await fetchSection('comparison-ajax', period);
                    statCards.value = data.statCards || [];
                    comparisonRows.value = data.comparisonRows || [];
                    previousPeriodLabel.value = data.previousPeriodLabel || '';
                } catch (e) {
                    comparisonError.value = e.message;
                } finally {
                    loadingComparison.value = false;
                }
            }

            async function fetchDashboardData(period) {
                loadingDaily.value = true;
                loadingTraffic.value = true;
                loadingHeatmap.value = true;
                dailyError.value = '';
                trafficError.value = '';
                heatmapError.value = '';
                try {
                    const data = await fetchSection('dashboard-data-ajax', period);
                    dailyStats.value = data.dailyStats || [];
                    totals.value = data.totals || {};
                    await nextTick();
                    renderChart(data.chartLabels || [], data.chartUniqueVisits || [], data.chartTotalVisits || []);
                    trafficCards.value = (data.trafficCards || []).map(c => reactive({ ...c, activeTab: 0 }));
                    heatmap.value = data.heatmap || { days: [], slots: [], intensities: [] };
                } catch (e) {
                    dailyError.value = e.message;
                    trafficError.value = e.message;
                    heatmapError.value = e.message;
                } finally {
                    loadingDaily.value = false;
                    loadingTraffic.value = false;
                    loadingHeatmap.value = false;
                }
            }

            function fetchAll(period) {
                // Fire independent requests in parallel
                fetchOverview(period);
                fetchComparison(period);
                fetchDashboardData(period);
            }

            function renderChart(labels, uniqueVisits, totalVisits) {
                if (!chartCanvas.value || !window.Chart) return;
                if (chartInstance) {
                    chartInstance.destroy();
                    chartInstance = null;
                }
                const ctx = chartCanvas.value.getContext('2d');
                chartInstance = new Chart(ctx, {
                    type: chartType.value,
                    data: {
                        labels: labels,
                        datasets: [
                            {
                                label: 'Unique Visitors',
                                backgroundColor: chartType.value === 'line' ? 'rgb(59, 130, 246)' : 'rgba(59, 130, 246, 0.5)',
                                borderColor: 'rgb(59, 130, 246)',
                                borderWidth: 2,
                                borderRadius: 4,
                                data: uniqueVisits,
                                pointRadius: chartType.value === 'line' ? 4 : 0,
                                tension: 0.2
                            },
                            {
                                label: 'Total Visitors',
                                backgroundColor: chartType.value === 'line' ? 'rgb(16, 185, 129)' : 'rgba(16, 185, 129, 0.5)',
                                borderColor: 'rgb(16, 185, 129)',
                                borderWidth: 2,
                                borderRadius: 4,
                                data: totalVisits,
                                pointRadius: chartType.value === 'line' ? 4 : 0,
                                tension: 0.2
                            }
                        ]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: { position: 'top', labels: { usePointStyle: true, padding: 20 } },
                            tooltip: { mode: 'index', intersect: false, padding: 10, bodySpacing: 5, backgroundColor: 'rgba(0, 0, 0, 0.8)' }
                        },
                        scales: {
                            y: { beginAtZero: true, grid: { drawBorder: false } },
                            x: { grid: { display: false } }
                        }
                    }
                });
            }

            function toggleChartType() {
                chartType.value = chartType.value === 'bar' ? 'line' : 'bar';
                const btn = document.getElementById('toggleChartType');
                if (btn) {
                    btn.innerHTML = chartType.value === 'line'
                        ? '<i class="bi bi-bar-chart"></i> Switch to Bar'
                        : '<i class="bi bi-graph-up"></i> Switch to Line';
                }
                if (chartInstance) {
                    const labels = chartInstance.data.labels;
                    const uniqueVisits = chartInstance.data.datasets[0].data;
                    const totalVisits = chartInstance.data.datasets[1].data;
                    renderChart(labels, uniqueVisits, totalVisits);
                }
            }

            function onPeriodChange() {
                fetchAll(selectedPeriod.value);
            }

            function heatmapColor(level) {
                const colors = ['#1c2333', '#12e198', '#14c987', '#17b176', '#1a9a65', '#1f8254'];
                return colors[level] || colors[0];
            }

            onMounted(() => {
                const urlParams = new URLSearchParams(window.location.search);
                const period = urlParams.get('period') || 'this-week';
                selectedPeriod.value = period;

                if (!window.Chart) {
                    const script = document.createElement('script');
                    script.src = 'https://cdn.jsdelivr.net/npm/chart.js';
                    script.onload = () => fetchAll(period);
                    document.head.appendChild(script);
                } else {
                    fetchAll(period);
                }

                const liveInterval = setInterval(() => {
                    const formData = new FormData();
                    formData.set('action', 'live-ajax');
                    fetch(buildApiUrl(), { method: 'POST', body: formData })
                        .then(r => r.json())
                        .then(d => { if (d.status === 'success' && d.data && d.data.liveVisitorCount !== undefined) liveVisitorCount.value = d.data.liveVisitorCount; })
                        .catch(() => {});
                }, 30000);

                onUnmounted(() => clearInterval(liveInterval));
            });

            return {
                selectedPeriod, periodOptions, liveVisitorCount,
                statCards, comparisonRows, previousPeriodLabel, dailyStats, totals,
                trafficCards, heatmap, metrics, selectedMetric, chartCanvas,
                loadingOverview, loadingComparison, loadingDaily, loadingTraffic, loadingHeatmap,
                overviewError, comparisonError, dailyError, trafficError, heatmapError,
                exportUrl, visitorActivityUrl, visitorPathsUrl,
                onPeriodChange, toggleChartType, heatmapColor
            };
        }
    }).mount('#dashboard-app');
})();
