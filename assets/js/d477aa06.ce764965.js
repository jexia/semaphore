(self.webpackChunksemaphore=self.webpackChunksemaphore||[]).push([[649],{3905:function(e,t,n){"use strict";n.d(t,{Zo:function(){return l},kt:function(){return d}});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function c(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?c(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):c(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function a(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},c=Object.keys(e);for(r=0;r<c.length;r++)n=c[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var c=Object.getOwnPropertySymbols(e);for(r=0;r<c.length;r++)n=c[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var u=r.createContext({}),s=function(e){var t=r.useContext(u),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},l=function(e){var t=s(e.components);return r.createElement(u.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},f=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,c=e.originalType,u=e.parentName,l=a(e,["components","mdxType","originalType","parentName"]),f=s(n),d=o,m=f["".concat(u,".").concat(d)]||f[d]||p[d]||c;return n?r.createElement(m,i(i({ref:t},l),{},{components:n})):r.createElement(m,i({ref:t},l))}));function d(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var c=n.length,i=new Array(c);i[0]=f;var a={};for(var u in t)hasOwnProperty.call(t,u)&&(a[u]=t[u]);a.originalType=e,a.mdxType="string"==typeof e?e:o,i[1]=a;for(var s=2;s<c;s++)i[s]=n[s];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}f.displayName="MDXCreateElement"},1572:function(e,t,n){"use strict";n.r(t),n.d(t,{frontMatter:function(){return i},metadata:function(){return a},toc:function(){return u},default:function(){return l}});var r=n(2122),o=n(9756),c=(n(7294),n(3905)),i={id:"functions",title:"Functions",sidebar_label:"Functions",slug:"/functions"},a={unversionedId:"functions",id:"functions",isDocsHomePage:!1,title:"Functions",description:"Functions could be used to preform computation on properties during runtime. Functions have read access to the entire reference store but could only write to their own stack.",source:"@site/docs/functions.md",sourceDirName:".",slug:"/functions",permalink:"/semaphore/docs/functions",editUrl:"https://github.com/jexia/semaphore/edit/master/website/docs/functions.md",version:"current",sidebar_label:"Functions",frontMatter:{id:"functions",title:"Functions",sidebar_label:"Functions",slug:"/functions"},sidebar:"docs",previous:{title:"Rollbacks",permalink:"/semaphore/docs/flows/rollbacks"},next:{title:"DevOps",permalink:"/semaphore/docs/devops"}},u=[],s={toc:u};function l(e){var t=e.components,n=(0,o.Z)(e,["components"]);return(0,c.kt)("wrapper",(0,r.Z)({},s,n,{components:t,mdxType:"MDXLayout"}),(0,c.kt)("p",null,"Functions could be used to preform computation on properties during runtime. Functions have read access to the entire reference store but could only write to their own stack.\nA unique resource is created for each function call where all references stored during runtime are located. This resource is created during compile time and references made to the given function are automatically adjusted."),(0,c.kt)("p",null,"A function should always return a property where all paths are absolute. This way it is easier for other properties to reference a resource."),(0,c.kt)("pre",null,(0,c.kt)("code",{parentName:"pre"},"function(...<arguments>)\n")),(0,c.kt)("p",null,"Functions could be called inside templates and could accept arguments and return a property as a response.\nA collection of predefined functions is included inside the Semaphore CLI."),(0,c.kt)("pre",null,(0,c.kt)("code",{parentName:"pre",className:"language-hcl"},'resource "auth" {\n    request "com.project" "Authenticate" {\n        header {\n            Authorization = "{{ jwt(input.header:Authorization) }}"\n        }\n    }\n}\n')))}l.isMDXComponent=!0}}]);