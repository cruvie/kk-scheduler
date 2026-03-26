package common_go

import (
	"github.com/cruvie/kk-scheduler/common_pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func MethodDescGetInterceptorAuth(method protoreflect.MethodDescriptor) common_pb.InterceptorAuth {
	options := method.Options()
	if options == nil {
		return common_pb.InterceptorAuth_UNSPECIFIED
	}

	interceptorList := proto.GetExtension(options, common_pb.E_InterceptorAuthList)
	if interceptorList == nil {
		return common_pb.InterceptorAuth_UNSPECIFIED
	}

	if list, ok := interceptorList.([]common_pb.InterceptorAuth); ok {
		switch len(list) {
		case 0:
			return common_pb.InterceptorAuth_UNSPECIFIED
		case 1:
			return list[0]
		}
		if len(list) > 1 {
			// only support one InterceptorAuth
			panic("MethodDescGetInterceptorAuth: len(list) > 1")
		}
	}

	return common_pb.InterceptorAuth_UNSPECIFIED
}

func MethodDescGetApiName(method protoreflect.MethodDescriptor) string {
	options := method.Options()
	if options == nil {
		return ""
	}

	serviceName := proto.GetExtension(options, common_pb.E_ApiName)
	if serviceName == nil {
		return ""
	}

	if name, ok := serviceName.(string); ok {
		return name
	}

	return ""
}
