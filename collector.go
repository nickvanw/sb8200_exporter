package sb8200exporter

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "sb8200"

var (
	labelsDownstream = []string{"downstream"}
	labelsUpstream   = []string{"upstream"}
)

// ModemCollector represents a single group of stats
type ModemCollector struct {
	c *client

	// Downstream metrics
	DownstreamFreq *prometheus.GaugeVec

	DownstreamPower *prometheus.GaugeVec

	DownstreamSNR *prometheus.GaugeVec

	DownstreamModulation *prometheus.GaugeVec

	DownstreamCorrecteds *prometheus.GaugeVec

	DownstreamUncorrectables *prometheus.GaugeVec

	// Upstrem Metrics

	UpstreamFreq *prometheus.GaugeVec

	UpstreamPower *prometheus.GaugeVec

	UpstreamWidth *prometheus.GaugeVec
}

// NewModemCollector creates a new statistics collector
func NewModemCollector(c *client) *ModemCollector {
	return &ModemCollector{
		c: c,

		DownstreamFreq: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstrem_freq_hertz",
				Help:      "Modem Downstream Frequency (Hz)",
			},
			labelsDownstream,
		),

		DownstreamPower: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_power_dbmv",
				Help:      "Modem Downstream Power (dBmV)",
			},
			labelsDownstream,
		),

		DownstreamSNR: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_snr_db",
				Help:      "Modem Downstream SNR (dB)",
			},
			labelsDownstream,
		),

		DownstreamModulation: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_modulation_qam",
				Help:      "Modem Downstream Modulation (QAM)",
			},
			labelsDownstream,
		),

		DownstreamCorrecteds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "dowmstream_correcteds_total",
				Help:      "Modem Downstream Correcteds",
			},
			labelsDownstream,
		),

		DownstreamUncorrectables: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "downstream_uncorrectables_total",
				Help:      "Modem Downstream Uncorrectables",
			},
			labelsDownstream,
		),

		UpstreamFreq: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "upstream_freq_hertz",
				Help:      "Modem Upstream Frequency (Hz)",
			},
			labelsUpstream,
		),

		UpstreamPower: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "upstream_power_dbmv",
				Help:      "Modem Upstream Power (dBmV)",
			},
			labelsUpstream,
		),

		UpstreamWidth: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "upstream_width_hertz",
				Help:      "Modem Upstream Channel Width (Hz)",
			},
			labelsUpstream,
		),
	}
}

func (m *ModemCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		m.DownstreamCorrecteds,
		m.DownstreamSNR,
		m.DownstreamModulation,
		m.DownstreamUncorrectables,
		m.DownstreamFreq,
		m.DownstreamPower,
		m.UpstreamFreq,
		m.UpstreamWidth,
		m.UpstreamPower,
	}
}

func (m *ModemCollector) collect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	data, err := m.c.fetch(ctx)
	if err != nil {
		return err
	}
	for _, node := range data.ds {
		downstreamID := strconv.Itoa(node.Channel)
		m.DownstreamSNR.WithLabelValues(downstreamID).Set(node.SNR)
		m.DownstreamFreq.WithLabelValues(downstreamID).Set(node.Freq)
		m.DownstreamPower.WithLabelValues(downstreamID).Set(node.Power)
		m.DownstreamModulation.WithLabelValues(downstreamID).Set(float64(node.Modulation))
		m.DownstreamCorrecteds.WithLabelValues(downstreamID).Set(float64(node.Correcteds))
		m.DownstreamUncorrectables.WithLabelValues(downstreamID).Set(float64(node.Uncorrectables))
	}

	for _, node := range data.us {
		upstreamID := strconv.Itoa(node.Channel)

		m.UpstreamFreq.WithLabelValues(upstreamID).Set(float64(node.Freq))
		m.UpstreamPower.WithLabelValues(upstreamID).Set(node.Power)
		m.UpstreamWidth.WithLabelValues(upstreamID).Set(float64(node.Width))
	}

	return nil
}

func (m *ModemCollector) describe(ch chan<- *prometheus.Desc) {
	for _, metric := range m.collectorList() {
		metric.Describe(ch)
	}
}
