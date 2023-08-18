// Code generated by scripts/content-script DO NOT EDIT.

package proxy

import "text/template"

const contentScript = "/* content-script v1.0.0 Sun Sep 08 2019 */\n" +
	"(function (configuration) {\n" +
	"    'use strict';\n" +
	"\n" +
	"    function getCurrentScript() {\n" +
	"        let { currentScript } = document;\n" +
	"        if (!currentScript) {\n" +
	"            const scripts = document.getElementsByTagName('script');\n" +
	"            currentScript = scripts[scripts.length - 1];\n" +
	"        }\n" +
	"        return currentScript;\n" +
	"    }\n" +
	"\n" +
	"    /**\n" +
	"     * Cosmetic rules object.\n" +
	"     *\n" +
	"     * @typedef {Object} Cosmeticresult\n" +
	"     * @property {StylesResult} elementHiding - Element hiding rules\n" +
	"     * @property {StylesResult} css - CSS rules\n" +
	"     * @property {ScriptsResult} js - JS rules\n" +
	"     */\n" +
	"\n" +
	"    /**\n" +
	"     * Styles result object\n" +
	"     *\n" +
	"     * @typedef {Object} StylesResult\n" +
	"     * @property {Array<string>} generic - Generic styles\n" +
	"     * @property {Array<string>} specific - Styles specific to this website\n" +
	"     * @property {Array<string>} genericExtCss - Generic ExtCSS styles\n" +
	"     * @property {Array<string>} specificExtCss - ExtCSS styles specific to this website\n" +
	"     */\n" +
	"\n" +
	"    /**\n" +
	"     * Scripts result object\n" +
	"     *\n" +
	"     * @typedef {Object} StylesResult\n" +
	"     * @property {Array<function>} generic - Generic functions\n" +
	"     * @property {Array<function>} specific - Functions specific to this website\n" +
	"     */\n" +
	"\n" +
	"    /**\n" +
	"     * Creates CSS rules from the cosmetic result\n" +
	"     * @param {*} rules rules\n" +
	"     * @param {*} style (optional) CSS style. For instance, `display: none`.\n" +
	"     */\n" +
	"    function getCssRules(rules, style) {\n" +
	"        const cssRules = [];\n" +
	"\n" +
	"        rules.forEach((rule) => {\n" +
	"            if (style) {\n" +
	"                cssRules.push(`${rule} { ${style} }`);\n" +
	"            } else {\n" +
	"                cssRules.push(rule);\n" +
	"            }\n" +
	"        });\n" +
	"\n" +
	"        return cssRules;\n" +
	"    }\n" +
	"\n" +
	"    /**\n" +
	"     * Creates a <style> tag that will be added to this page\n" +
	"     * @param {string} nonce - nonce string (that is added to the CSP of this page)\n" +
	"     * @param {CosmeticResult} cosmeticResult - cosmetic rules\n" +
	"     */\n" +
	"    function createStyle(nonce, cosmeticResult) {\n" +
	"        const style = document.createElement('style');\n" +
	"        style.setAttribute('nonce', nonce);\n" +
	"        style.setAttribute('type', 'text/css');\n" +
	"\n" +
	"        const cssRules = [\n" +
	"            ...getCssRules(cosmeticResult.elementHiding.generic, 'display: none!important'),\n" +
	"            ...getCssRules(cosmeticResult.elementHiding.specific, 'display: none!important'),\n" +
	"            ...getCssRules(cosmeticResult.css.generic),\n" +
	"            ...getCssRules(cosmeticResult.css.specific),\n" +
	"        ];\n" +
	"\n" +
	"        const cssTextNode = document.createTextNode(cssRules.join('\\n'));\n" +
	"        style.appendChild(cssTextNode);\n" +
	"        return style;\n" +
	"    }\n" +
	"\n" +
	"    /**\n" +
	"     * Applies cosmetic rules to the page\n" +
	"     *\n" +
	"     * @param {string} nonce - nonce string (that is added to the CSP of this page)\n" +
	"     * @param {CosmeticResult} cosmeticResult - cosmetic rules\n" +
	"     */\n" +
	"    function applyCosmeticResult(nonce, cosmeticResult) {\n" +
	"        const style = createStyle(nonce, cosmeticResult);\n" +
	"\n" +
	"        const currentScript = getCurrentScript();\n" +
	"        const rootElement = currentScript.parentNode;\n" +
	"        let insertBeforeElement = currentScript;\n" +
	"        if (currentScript.parentNode !== rootElement) {\n" +
	"            insertBeforeElement = null;\n" +
	"        }\n" +
	"        rootElement.insertBefore(style, insertBeforeElement);\n" +
	"\n" +
	"        /* Override styleEl's disabled\" property for forever enabled */\n" +
	"        const disabledDescriptor = {\n" +
	"            get: () => false,\n" +
	"            set: () => false,\n" +
	"        };\n" +
	"        Object.defineProperty(style, 'disabled', disabledDescriptor);\n" +
	"        Object.defineProperty(style.sheet, 'disabled', disabledDescriptor);\n" +
	"    }\n" +
	"\n" +
	"    // eslint-disable-next-line import/no-unresolved\n" +
	"\n" +
	"    const contentScriptExecutionFlagToCheck = configuration.nonce || 'adgRunId';\n" +
	"    if (!document[contentScriptExecutionFlagToCheck]) {\n" +
	"        // content script was already executed, doing nothing\n" +
	"        document[contentScriptExecutionFlagToCheck] = true;\n" +
	"        applyCosmeticResult(configuration.nonce, configuration.cosmeticResult);\n" +
	"    }\n" +
	"\n" +
	"}({\n" +
	"    \"nonce\": \"{{.Nonce}}\",\n" +
	"    \"cosmeticResult\": {\n" +
	"        \"elementHiding\": {\n" +
	"            \"generic\": [\n" +
	"                {{range .Result.ElementHiding.Generic}}\"{{js .}}\",{{end}}\n" +
	"            ],\n" +
	"            \"specific\": [\n" +
	"                {{range .Result.ElementHiding.Specific}}\"{{js .}}\",{{end}}\n" +
	"            ],\n" +
	"            \"genericExtCss\": [\n" +
	"                {{range .Result.ElementHiding.GenericExtCSS}}\"{{js .}}\",{{end}}\n" +
	"            ],\n" +
	"            \"specificExtCss\": [\n" +
	"                {{range .Result.ElementHiding.SpecificExtCSS}}\"{{js .}}\",{{end}}\n" +
	"            ],\n" +
	"        },\n" +
	"        \"css\": {\n" +
	"            \"generic\": [\n" +
	"                {{range .Result.CSS.Generic}}\"{{js .}}\",{{end}}\n" +
	"            ],\n" +
	"            \"specific\": [\n" +
	"                {{range .Result.CSS.Specific}}\"{{js .}}\",{{end}}\n" +
	"            ],\n" +
	"            \"genericExtCss\": [\n" +
	"                {{range .Result.CSS.GenericExtCSS}}\"{{js .}}\",{{end}}\n" +
	"            ],\n" +
	"            \"specificExtCss\": [\n" +
	"                {{range .Result.CSS.SpecificExtCSS}}\"{{js .}}\",{{end}}\n" +
	"            ],\n" +
	"        },\n" +
	"        \"js\": {\n" +
	"            \"generic\": [\n" +
	"                {{range .Result.JS.Generic}}() => { {{.}} },{{end}}\n" +
	"            ],\n" +
	"            \"specific\": [\n" +
	"                {{range .Result.JS.Specific}}() => { {{.}} },{{end}}\n" +
	"            ],\n" +
	"        }\n" +
	"    }\n" +
	"}));\n" +
	"\n"

var contentScriptTmpl = template.Must(template.New("contentScript").Parse(contentScript))
