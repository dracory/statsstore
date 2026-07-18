(function() {
    const { createApp, ref, reactive, computed, onMounted, watch, nextTick } = Vue;

    createApp({
        setup() {
            const loading = ref(true);
            const error = ref('');
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

            const exportUrl = computed(() => {
                return window.location.pathname + '?path=/admin/home&action=export&period=' + selectedPeriod.value;
            });

            function buildApiUrl(period) {
                const params = new URLSearchParams();
                params.set('path', '/admin/home');
                params.set('action', 'json');
                if (period) params.set('period', period);
                return window.location.pathname + '?' + params.toString();
            }

            async function fetchData(period) {
                loading.value = true;
                error.value = '';
                try {
                    const resp = await fetch(buildApiUrl(period));
                    if (!resp.ok) throw new Error('HTTP ' + resp.status);
                    const data = await resp.json();
                    if (data.error) throw new Error(data.error);

                    periodOptions.value = data.periodOptions || [];
                    selectedPeriod.value = data.selectedPeriod || 'this-week';
                    liveVisitorCount.value = data.liveVisitorCount || 0;
                    statCards.value = data.statCards || [];
                    comparisonRows.value = data.comparisonRows || [];
                    previousPeriodLabel.value = data.previousPeriodLabel || '';
                    dailyStats.value = data.dailyStats || [];
                    totals.value = data.totals || {};
                    trafficCards.value = (data.trafficCards || []).map(c => reactive({ ...c, activeTab: 0 }));
                    heatmap.value = data.heatmap || { days: [], slots: [], intensities: [] };

                    await nextTick();
                    renderChart(data.chartLabels || [], data.chartUniqueVisits || [], data.chartTotalVisits || []);
                } catch (e) {
                    error.value = e.message;
                } finally {
                    loading.value = false;
                }
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
                fetchData(selectedPeriod.value);
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
                    script.onload = () => fetchData(period);
                    document.head.appendChild(script);
                } else {
                    fetchData(period);
                }

                setInterval(() => {
                    fetch(buildApiUrl(selectedPeriod.value).replace('action=json', 'action=live'))
                        .then(r => r.json())
                        .then(d => { if (d.liveVisitorCount !== undefined) liveVisitorCount.value = d.liveVisitorCount; })
                        .catch(() => {});
                }, 30000);
            });

            return {
                loading, error, selectedPeriod, periodOptions, liveVisitorCount,
                statCards, comparisonRows, previousPeriodLabel, dailyStats, totals,
                trafficCards, heatmap, metrics, selectedMetric, chartCanvas,
                exportUrl, onPeriodChange, toggleChartType, heatmapColor
            };
        }
    }).mount('#dashboard-app');
})();
