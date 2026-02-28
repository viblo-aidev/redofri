package ixbrl

// writeCSS writes the inline CSS stylesheet for the annual report.
// Extracted and cleaned from the reference example file.
func (g *generator) writeCSS() {
	g.raw(`
html {
	font: 100%/1.4 Frutiger, Frutiger Linotype, Univers, DejaVu Sans Condensed, Liberation Sans, Nimbus Sans L, Geneva, Helvetica Neue, Helvetica, Arial, Tahoma, sans-serif; font-size-adjust: none; font-stretch: normal;
}
p {
	margin: 0px 0px 1.4em;
}
table {
	margin: 0px 0px 1.4em;
}
em {
	letter-spacing: -1px; font-style: normal; font-weight: bolder;
}
h1 {
	margin: 0px 0px 0.5em; line-height: 1; font-size: 180%;
}
h2 {
	margin: 2em 0px 0.5em; line-height: 1; font-size: 120%; font-weight: bold;
}
h1 + h2 {
	margin-top: 0px;
}
h3 {
	margin: 1em 0px 0.5em; line-height: 1; font-size: 110%; font-weight: normal;
}
h4 {
	margin: 1.3em 0px 0.5em; line-height: 1; font-size: 100%; font-weight: normal;
}
h5 {
	margin: 1.61em 0px 0.5em; line-height: 1; font-size: 90%; font-style: italic; font-weight: normal;
}
h6 {
	margin: 2em 0px 0.5em; line-height: 1; font-size: 80%; font-weight: normal;
}
img {
	height: auto; max-width: 100%;
}
table {
	width: 100%; border-collapse: collapse;
}
td {
	padding: 0px 0.5em; text-align: left; font-size: smaller; vertical-align: top;
}
th {
	padding: 0px 0.5em; text-align: left; font-size: smaller; vertical-align: top;
}
tbody th[colspan] {
	border-width: 1px 0px; border-style: solid; border-color: rgb(0, 0, 0); border-image: none; font-weight: bold;
}
tbody th[colspan='1'] {
	border: 0px currentColor; border-image: none; font-weight: normal;
}
@media screen {
body {
	margin: 0px auto; padding: 0px; max-width: 70em;
}
h1, h2, h3, h4, h5, h6 {
	color: rgb(68, 68, 68);
}
#wrapper {
	padding: 0px 0px 0px 15em; background-color: rgb(255, 255, 255);
}
}
@media print {
h1, h2, h3, h4, h5, h6 {
	color: rgb(0, 0, 0);
}
table {
	page-break-inside: avoid;
}
td, th {
	border: 1px solid rgb(153, 153, 153); border-image: none;
}
thead td, thead th {
	border-width: 0px 0px 1px; border-color: rgb(0, 0, 0);
}
#wrapper {
	background: none; color: rgb(0, 0, 0);
}
abbr[title] {
	border-bottom-color: currentColor; border-bottom-width: 0px; border-bottom-style: none;
}
}
.ar-page {
	margin: 1em 0px; padding: 1em 2em 2em; border: 1px solid rgb(240, 240, 240); border-image: none; line-height: 1.2; min-height: 52em; max-width: 38em; box-shadow: 0.25em 0.25em 0.3em #999;
}
.note.ar-page {
	min-width: 20em;
}
.wide.ar-page {
	min-width: 30em;
}
.ar-page {
	font-family: Times New Roman,Times,serif;
}
#main .ar-page dl, #main .ar-page p {
	font-family: Times New Roman,Times,serif;
}
.ar-page abbr[title] {
	border: 0px currentColor; border-image: none;
}
.ar-page dt {
	text-decoration: underline;
}
.ar-page dd {
	margin: 0px;
}
.ar-page h3 {
	margin: 0px 0px 0.5em; font-size: 100%; font-weight: bold;
}
.ar-page h3[id] {
	line-height: 1.8em;
}
.ar-page h4 {
	margin: 1.5em 0px; font-size: 100%; font-style: italic; font-weight: normal;
}
.ar-page h3 + h4 {
	margin: 0px 0px 0.5em;
}
.ar-page .join {
	margin-bottom: 0px;
}
.ar-page span.note {
	margin-right: 1.5em;
}
.ar-page table tbody tr {
	background-color: transparent;
}
.ar-page td, .ar-page th {
	border: 0px currentColor; border-image: none; padding-left: 0px; font-size: 1em; vertical-align: bottom;
}
.ar-page td[rowspan] {
	vertical-align: top;
}
.ar-page th {
	text-align: left; font-weight: normal;
}
.ar-page th[colspan] {
	border: 0px currentColor; border-image: none;
}
table.col-3 td + td, table.col-3 th + th {
	text-align: right;
}
table.col-4 {
	min-width: 30em;
}
table.col-4 td + td + td, table.col-4 th + th + th {
	text-align: right;
}
table.col-5 td + td, table.col-5 th + th {
	text-align: right;
}
.ar-financial col.kr {
	width: 8em;
}
.ar-financial thead th {
	font-size: 120%; font-weight: bold;
}
.ar-equity.ar-financial thead th {
	font-size: 100%; font-weight: normal;
}
.ar-financial thead th[colspan] {
	padding-top: 1em; padding-bottom: 1em;
}
.ar-financial tbody tr:first-child th {
	padding-top: 1em; font-weight: bold;
}
.ar-financial tbody tr:first-child td {
	padding-top: 1em;
}
.ar-financial tbody th.sub {
	font-style: italic; font-weight: normal;
}
.ar-financial tbody th.sup {
	font-size: 120%; font-weight: bold;
}
.ar-financial tbody .sep {
	padding-top: 1em;
}
.ar-financial tbody tr.sep td, .ar-financial tbody tr.sep th {
	padding-top: 1em;
}
.ar-financial tr.result td {
	padding-top: 1em;
}
.ar-financial tr.result td:first-child {
	font-style: italic; font-weight: bold;
}
.ar-financial td a {
	color: rgb(0, 0, 0); text-decoration: none;
}
table.ar-note {
	table-layout: fixed; min-width: 20em;
}
.ar-note col.kr {
	width: 8em;
}
.ar-note td + td, .ar-note th + th, .ar-note th + td {
	text-align: right;
}
.ar-note thead th span {
	border-bottom-color: currentColor; border-bottom-width: 1px; border-bottom-style: solid;
}
.ar-note tbody tr:first-child th {
	padding-top: 1em; font-weight: bold;
}
.ar-note tbody tr:first-child td {
	padding-top: 1em;
}
table.ar-note + h3 {
	margin-top: 3em;
}
table.ar-note-10 {
	table-layout: fixed; min-width: 20em; max-width: 30em;
}
.ar-note-10 col.kr {
	width: 8em;
}
.ar-note-10 td + td {
	text-align: right;
}
col.kr {
	width: 6em;
}
col.note {
	width: 3em;
}
col.tkr {
	width: 5em;
}
td.sub-sum {
	font-style: italic;
}
td.sum {
	font-weight: bold;
}
td.total {
	font-size: 120%;
}
td .sum {
	border-bottom-color: currentColor; border-bottom-width: 1px; border-bottom-style: solid;
}
td .total {
	border-bottom-color: currentColor; border-bottom-width: 3px; border-bottom-style: double;
}
.ar-page-hdr {
	text-align: right; color: rgb(102, 102, 102); overflow: hidden; margin-bottom: 3em; white-space: nowrap; font-size: 75%;
}
.ar-page-hdr span {
	text-align: left; float: left;
}
.ar-toc th {
	padding-bottom: 1em;
}
.ar-toc td + td, .ar-toc th + th {
	text-align: right;
}
.ar-toc td a {
	color: rgb(0, 0, 0); text-decoration: none;
}
.ar-toc td span {
	margin-right: 1em;
}
#main #ar-certification > p {
	font-family: Calibri,Tahoma,sans-serif;
}
.ar-logo {
	margin: 5em 0px 4em; font-weight: bold;
}
.ar-signature {
	font-family: Calibri,Tahoma,sans-serif; margin-left: 0.5em;
}
.ar-signature-text {
	padding-left: 0.5em; font-weight: bold; margin-top: 1em;
}
.ar-signature-label {
	padding-left: 0.6em; font-size: smaller; border-top-color: currentColor; border-top-width: 1px; border-top-style: dotted;
}
.ar-overview {
	table-layout: fixed;
}
.ar-overview thead th {
	font-size: 100%; font-weight: normal;
}
.ar-overview tbody tr:first-child td {
	padding-top: 0px;
}
p.ar-disp {
	margin: 1.5em 0px 0px;
}
table.ar-disp {
	width: 75%; table-layout: fixed; min-width: 16em;
}
.ar-disp td + td {
	text-align: right;
}
.ar-dividend td {
	padding-top: 1em;
}
#ar3-page-3 h3 {
	margin: 0px; font-weight: normal; text-decoration: underline;
}
.ar-profit-loss {
	table-layout: fixed;
}
#ar3-page-4 .ar-financial tbody:last-child tr.result td {
	padding-top: 0px;
}
.ar-balance-sheet {
	table-layout: fixed;
}
.ar-depreciation {
	width: auto;
}
.ar-depreciation td:first-child {
	min-width: 15em;
}
.ar-depreciation td + td {
	text-align: right;
}
.ar-capital {
	table-layout: fixed;
}
.ar-capital td + td, .ar-capital th + th {
	text-align: right;
}
.ar-capital th {
	font-weight: normal; vertical-align: top;
}
.ar-capital thead th {
	padding-right: 0px; padding-left: 0.5em;
}
.ar-capital td + td {
	padding-right: 0px; padding-left: 0.5em;
}
col.kr-2 {
	width: 4.5em;
}
.ar-capital + h3 {
	margin-top: 3em;
}
.ar-signature-2 {
	margin-top: 2em;
}
.ar-signature-2 div.name {
	width: 40%; padding-top: 0em; vertical-align: top; display: inline-block;
}
@media print {
.ar-page {
	margin: 0px; border: 0px currentColor; border-image: none; font-size: 11pt; page-break-after: always; min-height: 0px; max-width: none; page-break-inside: avoid; box-shadow: none;
}
.ar-page:last-of-type {
	page-break-after: avoid;
}
}
`)
}
