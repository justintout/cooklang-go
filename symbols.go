package cooklang

const eof = '\u0000'

const leftEscape = `\`

const leftLineComment = "--"
const leftBlockComment = "[-"
const rightBlockComment = "-]"

const leftMetadata = ">>"

const leftIngredient = "@"
const leftCookware = "#"
const leftTimer = "~"

const leftQuantity = "{"
const rightQuantity = "}"
const dividerQuantity = "%"
