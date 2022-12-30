package hitcounter

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type HitCounterProps struct {
	Downstream awslambda.IFunction
}

type hitCounter struct {
	constructs.Construct
	handler awslambda.IFunction
	table   awsdynamodb.Table
}

type HitCounter interface {
	constructs.Construct
	Handler() awslambda.IFunction
	Table() awsdynamodb.Table
}

func NewHitCounter(scope constructs.Construct, id string, props *HitCounterProps) HitCounter {
	this := constructs.NewConstruct(scope, &id)

	table := awsdynamodb.NewTable(this, jsii.String("lggocdkHits"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("path"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	handler := awslambda.NewFunction(this, jsii.String("lggocdkHitCounterHandler"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_16_X(),
		Handler: jsii.String("hitcounter.handler"),
		Code:    awslambda.Code_FromAsset(jsii.String("lambda"), nil),
		Environment: &map[string]*string{
			"HITS_TABLE_NAME":          table.TableName(),
			"DOWNSTREAM_FUNCTION_NAME": props.Downstream.FunctionName(),
		},
	})

	table.GrantReadWriteData(handler)
	props.Downstream.GrantInvoke(handler)

	return &hitCounter{this, handler, table}
}

func (h *hitCounter) Handler() awslambda.IFunction {
	return h.handler
}

func (h *hitCounter) Table() awsdynamodb.Table {
	return h.table
}
