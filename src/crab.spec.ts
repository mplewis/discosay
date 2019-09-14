import crab from "./crab";

import { readFileSync } from 'fs-extra'
import { join } from "path";

const renderedCrab = readFileSync(join(__dirname, '..', 'spec', 'fixtures', 'crab.txt')).toString().trim()

describe("crab", () => {
  it("renders as expected", () => {
    expect(crab("You should be using Rust!")).toMatch(renderedCrab);
  });
});
