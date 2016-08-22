"use strict";

var assert = require("chai").assert;
var codecheck = require("codecheck");
var TestcaseReader = require("testcase-util");

var appName = process.env.APP_COMMAND

describe("いろは市", () => {
  var app = codecheck.consoleApp(appName);
  var reader = new TestcaseReader("./test/testcase1.md");

  reader.testcases().forEach(v => {
    it(v.title, () => {
      return app.codecheck(v.input).then(result => {
        assert.equal(result.code, 0);
        assert.deepEqual(result.stdout, v.output);
      });
    });
  });
});
