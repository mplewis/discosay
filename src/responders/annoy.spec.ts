import annoy from "./annoy";

import fakeMessage from '../../spec/fake_message'

function applicableCount(message: string): number {
  const fake = fakeMessage(message);
  const success = [];
  for (let i = 0; i < 100; i++) {
    success.push(annoy.applicable(fake))
  }
  return success.filter(Boolean).length
}

describe('annoy', () => {
  it('applies to messages containing "rust"', () => {
    const messages = [
      "rust is great!",
      "i love Rust!",
      "did you write it in RUST?"
    ]
    for (let message of messages) {
      expect(applicableCount(message)).toBeGreaterThan(0)
    }
  })

  it('does not apply to messages containing "rust"', () => {
    const messages = [
      "C++ is a good language",
      "Python is nice because it has no types",
      "Most things should be rewritten in OCaml."
    ]
    for (let message of messages) {
      expect(applicableCount(message)).toEqual(0)
    }
  })
})
