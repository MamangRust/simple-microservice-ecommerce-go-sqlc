package helperproto

import "google.golang.org/protobuf/types/known/wrapperspb"

func StringPtrToWrapper(s *string) *wrapperspb.StringValue {
	if s == nil {
		return nil
	}
	return wrapperspb.String(*s)
}
