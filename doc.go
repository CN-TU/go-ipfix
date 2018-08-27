/*
Package ipfix writes ipfix data streams as defined by RFC 7011.

Currently supported is only writing to an io.Writer, the datatypes from the
RFC7011 + basic lists from RFC 6313.

Option templates and template recovation are not supported.

Usage

For exporting ipfix data a MessageStream instance has to be created with MakeMessageStream.
This stream then provides the two functions AddTemplate for adding templates and SendData for sending
data, as specified by a template. After all the data has been added with SendData, Flush must be called.
Full examples are provided at the MakeMessageStream function.

Information elements can be created either from an iespec (RFC 7373) with MakeIEFromSpec, or by hand
with NewInformationElement or NewBasicList.

All the information elements as defined by the iana can be loaded with LoadIanaSpec and then accessed
by name with GetInformationElement.

Custom loaders can be created with generate_spec.go or by calling RegisterInformationElement.

*/
package ipfix
