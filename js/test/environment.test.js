import { describe,expect,it } from "bun:test";
import process from "process";

describe("Testing Environment Variable", () => {
    const config = process.env
    const testCases = [
        {
            name : "should be exist key environment",
            value : config["APP_ENV"],
            expectValue : "testing",
        },
        {
            name : "should be exist key environment",
            value : config["SECRET_KEY"],
            expectValue : "secret",
        },
        {
            name : "should be exist key environment",
            value : config["MESSAGE_BROKER_PROTOCOL"],
            expectValue : "amqp",
        },
        {
            name : "should be exist key environment",
            value : config["MESSAGE_BROKER_HOST"],
            expectValue : "localhost",
        },
        {
            name : "should be exist key environment",
            value : config["MESSAGE_BROKER_PORT"],
            expectValue : "15672",
        },
        {
            name : "should be exist key environment",
            value : config["MESSAGE_BROKER_USER"],
            expectValue : "guest",
        },
        {
            name : "should be exist key environment",
            value : config["MESSAGE_BROKER_PASSWORD"],
            expectValue : "guest",
        },
        {
            name : "should be exist key environment",
            value : config["MESSAGE_BROKER_VIRTUAL_HOST"],
            expectValue : "testing",
        },
        {
            name : "should be not found exist key environment",
            value : config["APP_ENVS"],
            expectValue : undefined,
        },
        {
            name : "should be not found exist key environment",
            value : config["SECRET_KEYS"],
            expectValue : undefined,
        },
        {
            name : "should be not found exist key environment",
            value : config["MESSAGE_BROKER_PROTOCOLS"],
            expectValue : undefined,
        },
        {
            name : "should be not found exist key environment",
            value : config["MESSAGE_BROKER_HOSTS"],
            expectValue : undefined,
        },
        {
            name : "should be not found exist key environment",
            value : config["MESSAGE_BROKER_PORTS"],
            expectValue : undefined,
        },
        {
            name : "should be not found exist key environment",
            value : config["MESSAGE_BROKER_USERS"],
            expectValue : undefined,
        },
        {
            name : "should be not found exist key environment",
            value : config["MESSAGE_BROKER_PASSWORDS"],
            expectValue : undefined,
        },
        {
            name : "should be not found exist key environment",
            value : config["MESSAGE_BROKER_VIRTUAL_HOSTS"],
            expectValue : undefined,
        },
    ]
    testCases.forEach((testCase,index) => {
        index++
        it(`Case ${index} : ${testCase.name}`,() => {
            expect(typeof(testCase.expectValue)).toBe(typeof(testCase.value));
        })
    })
});