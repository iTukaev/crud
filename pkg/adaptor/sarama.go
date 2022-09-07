package adaptor

import "github.com/Shopify/sarama"

func ConsumerHeaderToProducer(cHeaders []*sarama.RecordHeader) []sarama.RecordHeader {
	pHeaders := make([]sarama.RecordHeader, len(cHeaders))
	for _, head := range cHeaders {
		pHeaders = append(pHeaders, *head)
	}
	return pHeaders
}
