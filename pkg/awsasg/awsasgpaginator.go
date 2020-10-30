package awsasg

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

// ASGDescribeAutoScalingGroupsAPIClient provides interface for the S3 API client
// ListObjectsV2 operation call.
type ASGDescribeAutoScalingGroupsAPIClient interface {
	DescribeAutoScalingGroups(context.Context, *autoscaling.DescribeAutoScalingGroupsInput, ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
}

// ASGDescribeAutoScalingGroupsPaginatorOptions provides the options for configuring the
// ListObjectsV2Paginator.
type ASGDescribeAutoScalingGroupsPaginatorOptions struct {
	// The maximum number of keys to return per page.
	Limit int32
}

// ASGDescribeAutoScalingGroupsPaginator provides the paginator to paginate S3 ListObjectsV2
// response pages.
type ASGDescribeAutoScalingGroupsPaginator struct {
	options ASGDescribeAutoScalingGroupsPaginatorOptions

	client ASGDescribeAutoScalingGroupsAPIClient
	params autoscaling.DescribeAutoScalingGroupsInput

	nextToken *string
	firstPage bool
}

// NewASGDescribeAutoScalingGroupsPaginator initializes a new S3 ListObjectsV2 Paginator for
// paginating the ListObjectsV2 respones.
func NewASGDescribeAutoScalingGroupsPaginator(client ASGDescribeAutoScalingGroupsAPIClient, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*ASGDescribeAutoScalingGroupsPaginatorOptions)) *ASGDescribeAutoScalingGroupsPaginator {
	var options ASGDescribeAutoScalingGroupsPaginatorOptions
	for _, fn := range optFns {
		fn(&options)
	}
	p := &ASGDescribeAutoScalingGroupsPaginator{
		options:   options,
		client:    client,
		firstPage: true,
	}
	if params != nil {
		p.params = *params
	}
	return p
}

// HasMorePages returns true if there are more pages or if the first page has
// not been retrieved yet.
func (p *ASGDescribeAutoScalingGroupsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)

}

// NextPage attempts to retrieve the next page, or returns error if unable to.
func (p *ASGDescribeAutoScalingGroupsPaginator) NextPage(ctx context.Context) (
	*autoscaling.DescribeAutoScalingGroupsOutput, error,
) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := p.params
	result, err := p.client.DescribeAutoScalingGroups(ctx, &params)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	if result.NextToken == nil {
		p.nextToken = nil
	} else {
		p.nextToken = result.NextToken
	}
	p.params.NextToken = p.nextToken

	return result, nil

}
