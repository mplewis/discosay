const template = `
      \\
       \\
         _~^~^~_
     \\) /  o o  \\ (/
       '_   -   _'
       / '-----' \\
`

function bubble (msg: string): string {
  const width = 4 + msg.length
  let topbot = ''
  for (let i = 0; i < width; i++) {
    topbot += '-'
  }
  const middle = `| ${msg} |`
  return [topbot, middle, topbot].join('\n')
}

export default function crab (msg: string): string {
  return '```\n' + bubble(msg) + template + '```'
}
