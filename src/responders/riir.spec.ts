import riir from "./riir";
import fakeMessage from '../../spec/fake_message'

describe('riir', () => {
  it('applies to messages matching the trigger', () => {
    expect(riir.applicable(fakeMessage('riir'))).toBe(true);
    expect(riir.applicable(fakeMessage('rewriir'))).toBe(false);
    expect(riir.applicable(fakeMessage('riir today'))).toBe(false);
  })
})
