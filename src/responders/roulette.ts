import { Responder } from '../types'
import crab from '../crab'

const trigger = '!retf'
const messages = [
  'Have you tried rewriting it in Rust?',
  "You sure that doesn't have a memory leak?",
  'Rust is the best language and C is the worst',
  '>implying garbage collectors are even real',
  'Rust Evangelism Task Force v0.0.1',
  "Valgrind? No thanks, I don't do drugs",
  'Imagine not ensuring memory safety at compile time.',
  'bro, you just coded cringe. you are going to leak pointers'
]

const roulette: Responder = {
  name: 'roulette',
  applicable: msg => msg.content.toLowerCase() === trigger,
  handle: _ => {
    const count = messages.length
    const i = Math.floor(Math.random() * count)
    const msg = messages[i]
    return crab(msg)
  }
}

export default roulette
