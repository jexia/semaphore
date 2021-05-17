(self.webpackChunksemaphore=self.webpackChunksemaphore||[]).push([[877],{3905:function(e,t,n){"use strict";n.d(t,{Zo:function(){return p},kt:function(){return d}});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var s=r.createContext({}),u=function(e){var t=r.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},p=function(e){var t=u(e.components);return r.createElement(s.Provider,{value:t},e.children)},c={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},m=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,i=e.originalType,s=e.parentName,p=l(e,["components","mdxType","originalType","parentName"]),m=u(n),d=o,f=m["".concat(s,".").concat(d)]||m[d]||c[d]||i;return n?r.createElement(f,a(a({ref:t},p),{},{components:n})):r.createElement(f,a({ref:t},p))}));function d(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=n.length,a=new Array(i);a[0]=m;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:o,a[1]=l;for(var u=2;u<i;u++)a[u]=n[u];return r.createElement.apply(null,a)}return r.createElement.apply(null,n)}m.displayName="MDXCreateElement"},4358:function(e,t,n){"use strict";n.r(t),n.d(t,{frontMatter:function(){return a},metadata:function(){return l},toc:function(){return s},default:function(){return p}});var r=n(2122),o=n(9756),i=(n(7294),n(3905)),a={id:"installation.contributing",title:"Contributing",sidebar_label:"Contributing",slug:"/installation/contributing"},l={unversionedId:"installation.contributing",id:"installation.contributing",isDocsHomePage:!1,title:"Contributing",description:"If you wish to work on Semaphore itself or any of its built-in systems, you'll",source:"@site/docs/installation-contributing.md",sourceDirName:".",slug:"/installation/contributing",permalink:"/semaphore/docs/installation/contributing",editUrl:"https://github.com/jexia/semaphore/edit/master/website/docs/installation-contributing.md",version:"current",sidebar_label:"Contributing",frontMatter:{id:"installation.contributing",title:"Contributing",sidebar_label:"Contributing",slug:"/installation/contributing"},sidebar:"docs",previous:{title:"Build from source",permalink:"/semaphore/docs/installation/source"},next:{title:"Getting started",permalink:"/semaphore/docs/flows"}},s=[],u={toc:s};function p(e){var t=e.components,n=(0,o.Z)(e,["components"]);return(0,i.kt)("wrapper",(0,r.Z)({},u,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("p",null,"If you wish to work on Semaphore itself or any of its built-in systems, you'll\nfirst need ",(0,i.kt)("a",{parentName:"p",href:"https://www.golang.org"},"Go")," installed on your machine. Go version\n1.13.7+ is ",(0,i.kt)("em",{parentName:"p"},"required"),"."),(0,i.kt)("p",null,"For local dev first make sure Go is properly installed, including setting up a\n",(0,i.kt)("a",{parentName:"p",href:"https://golang.org/doc/code.html#GOPATH"},"GOPATH"),". Ensure that ",(0,i.kt)("inlineCode",{parentName:"p"},"$GOPATH/bin")," is in\nyour path as some distributions bundle old version of build tools. Next, clone this\nrepository. Semaphore uses ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/golang/go/wiki/Modules"},"Go Modules"),",\nso it is recommended that you clone the repository ",(0,i.kt)("strong",{parentName:"p"},(0,i.kt)("em",{parentName:"strong"},"outside"))," of the GOPATH.\nYou can then download any required build tools by bootstrapping your environment:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-sh"},"$ make bootstrap\n...\n")),(0,i.kt)("p",null,"To compile a development version of Semaphore, run ",(0,i.kt)("inlineCode",{parentName:"p"},"make")," or ",(0,i.kt)("inlineCode",{parentName:"p"},"make dev"),". This will\nput the Semaphore binary in the ",(0,i.kt)("inlineCode",{parentName:"p"},"bin")," folders:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-sh"},"$ make dev\n...\n$ bin/semaphore\n...\n")),(0,i.kt)("p",null,"To run tests, type ",(0,i.kt)("inlineCode",{parentName:"p"},"make test"),". If\nthis exits with exit status 0, then everything is working!"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-sh"},"$ make test\n...\n")))}p.isMDXComponent=!0}}]);