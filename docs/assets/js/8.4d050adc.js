(window.webpackJsonp=window.webpackJsonp||[]).push([[8],{358:function(t,s,n){"use strict";n.r(s);var e=n(42),a=Object(e.a)({},(function(){var t=this,s=t.$createElement,n=t._self._c||s;return n("ContentSlotsDistributor",{attrs:{"slot-key":t.$parent.slotKey}},[n("h1",{attrs:{id:"introduction"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#introduction"}},[t._v("#")]),t._v(" Introduction")]),t._v(" "),n("p",[t._v("Semaphore is a tool to orchestrate your micro-service architecture. Requests could be manipulated passed branched to different services to be returned as a single output.\nYou could define request flows on top of your currently existing schema definitions.\nPlease check out the examples directory for more examples.")]),t._v(" "),n("div",{staticClass:"custom-block tip"},[n("p",{staticClass:"custom-block-title"},[t._v("TIP")]),t._v(" "),n("p",[t._v("In many of the available examples are protobuffers used. Semaphore currently supports protobuffers more official schema definitions such as Avro and XML will be added in the future")])]),t._v(" "),n("div",{staticClass:"language-hcl extra-class"},[n("pre",{pre:!0,attrs:{class:"language-hcl"}},[n("code",[t._v("endpoint "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"GetUser"')]),t._v(" "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"http"')]),t._v(" "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("{")]),t._v("\n    "),n("span",{pre:!0,attrs:{class:"token property"}},[t._v("endpoint")]),t._v(" "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("=")]),t._v(" "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"/user/:id"')]),t._v("\n    "),n("span",{pre:!0,attrs:{class:"token property"}},[t._v("method")]),t._v(" "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("=")]),t._v(" "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"GET"')]),t._v("\n"),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("}")]),t._v("\n\nflow "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"GetUser"')]),t._v(" "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("{")]),t._v("\n    input "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"proto.Query"')]),t._v(" "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("{")]),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("}")]),t._v("\n    \n    resource "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"user"')]),t._v(" "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("{")]),t._v("\n        request "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"proto.Users"')]),t._v(" "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"Get"')]),t._v(" "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("{")]),t._v("\n            "),n("span",{pre:!0,attrs:{class:"token property"}},[t._v("id")]),t._v(" "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("=")]),t._v(" "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"{{ input:id }}"')]),t._v("\n        "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("}")]),t._v("\n    "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("}")]),t._v("\n    \n    "),n("span",{pre:!0,attrs:{class:"token keyword"}},[t._v("output"),n("span",{pre:!0,attrs:{class:"token type variable"}},[t._v(' "proto.User" ')])]),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("{")]),t._v("\n        "),n("span",{pre:!0,attrs:{class:"token property"}},[t._v("name")]),t._v(" "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("=")]),t._v(" "),n("span",{pre:!0,attrs:{class:"token string"}},[t._v('"{{ user:name }}"')]),t._v("\n    "),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("}")]),t._v("\n"),n("span",{pre:!0,attrs:{class:"token punctuation"}},[t._v("}")]),t._v("\n")])])])])}),[],!1,null,null,null);s.default=a.exports}}]);