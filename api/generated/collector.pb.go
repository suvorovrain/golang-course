





package generated

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RepoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Owner         string                 `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	Repo          string                 `protobuf:"bytes,2,opt,name=repo,proto3" json:"repo,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RepoRequest) Reset() {
	*x = RepoRequest{}
	mi := &file_api_proto_collector_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RepoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepoRequest) ProtoMessage() {}

func (x *RepoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_collector_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}


func (*RepoRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_collector_proto_rawDescGZIP(), []int{0}
}

func (x *RepoRequest) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *RepoRequest) GetRepo() string {
	if x != nil {
		return x.Repo
	}
	return ""
}

type RepoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	FullName      string                 `protobuf:"bytes,2,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Url           string                 `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	Stars         int32                  `protobuf:"varint,5,opt,name=stars,proto3" json:"stars,omitempty"`
	Forks         int32                  `protobuf:"varint,6,opt,name=forks,proto3" json:"forks,omitempty"`
	Watchers      int32                  `protobuf:"varint,7,opt,name=watchers,proto3" json:"watchers,omitempty"`
	Language      string                 `protobuf:"bytes,8,opt,name=language,proto3" json:"language,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     string                 `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Error         string                 `protobuf:"bytes,11,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RepoResponse) Reset() {
	*x = RepoResponse{}
	mi := &file_api_proto_collector_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RepoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepoResponse) ProtoMessage() {}

func (x *RepoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_collector_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}


func (*RepoResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_collector_proto_rawDescGZIP(), []int{1}
}

func (x *RepoResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RepoResponse) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *RepoResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *RepoResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *RepoResponse) GetStars() int32 {
	if x != nil {
		return x.Stars
	}
	return 0
}

func (x *RepoResponse) GetForks() int32 {
	if x != nil {
		return x.Forks
	}
	return 0
}

func (x *RepoResponse) GetWatchers() int32 {
	if x != nil {
		return x.Watchers
	}
	return 0
}

func (x *RepoResponse) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *RepoResponse) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *RepoResponse) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

func (x *RepoResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_api_proto_collector_proto protoreflect.FileDescriptor

const file_api_proto_collector_proto_rawDesc = "" +
	"\n" +
	"\x19api/proto/collector.proto\x12\tcollector\"7\n" +
	"\vRepoRequest\x12\x14\n" +
	"\x05owner\x18\x01 \x01(\tR\x05owner\x12\x12\n" +
	"\x04repo\x18\x02 \x01(\tR\x04repo\"\xab\x02\n" +
	"\fRepoResponse\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x1b\n" +
	"\tfull_name\x18\x02 \x01(\tR\bfullName\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12\x10\n" +
	"\x03url\x18\x04 \x01(\tR\x03url\x12\x14\n" +
	"\x05stars\x18\x05 \x01(\x05R\x05stars\x12\x14\n" +
	"\x05forks\x18\x06 \x01(\x05R\x05forks\x12\x1a\n" +
	"\bwatchers\x18\a \x01(\x05R\bwatchers\x12\x1a\n" +
	"\blanguage\x18\b \x01(\tR\blanguage\x12\x1d\n" +
	"\n" +
	"created_at\x18\t \x01(\tR\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\n" +
	" \x01(\tR\tupdatedAt\x12\x14\n" +
	"\x05error\x18\v \x01(\tR\x05error2R\n" +
	"\x10CollectorService\x12>\n" +
	"\vGetRepoInfo\x12\x16.collector.RepoRequest\x1a\x17.collector.RepoResponseB,Z*github-info-system/api/generated;generatedb\x06proto3"

var (
	file_api_proto_collector_proto_rawDescOnce sync.Once
	file_api_proto_collector_proto_rawDescData []byte
)

func file_api_proto_collector_proto_rawDescGZIP() []byte {
	file_api_proto_collector_proto_rawDescOnce.Do(func() {
		file_api_proto_collector_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_proto_collector_proto_rawDesc), len(file_api_proto_collector_proto_rawDesc)))
	})
	return file_api_proto_collector_proto_rawDescData
}

var file_api_proto_collector_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_proto_collector_proto_goTypes = []any{
	(*RepoRequest)(nil),  
	(*RepoResponse)(nil), 
}
var file_api_proto_collector_proto_depIdxs = []int32{
	0, 
	1, 
	1, 
	0, 
	0, 
	0, 
	0, 
}

func init() { file_api_proto_collector_proto_init() }
func file_api_proto_collector_proto_init() {
	if File_api_proto_collector_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_proto_collector_proto_rawDesc), len(file_api_proto_collector_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_collector_proto_goTypes,
		DependencyIndexes: file_api_proto_collector_proto_depIdxs,
		MessageInfos:      file_api_proto_collector_proto_msgTypes,
	}.Build()
	File_api_proto_collector_proto = out.File
	file_api_proto_collector_proto_goTypes = nil
	file_api_proto_collector_proto_depIdxs = nil
}
