# bj4

[![GoDoc](https://godoc.org/github.com/rayark/go-bj4?status.svg)](https://godoc.org/github.com/rayark/go-bj4)
[![Build Status](https://travis-ci.org/rayark/go-bj4.svg?branch=master)](https://travis-ci.org/rayark/go-bj4)

One-task-at-a-time task scheduling.

## Migration Guide (V0 -> V1)
Here are the key changes introduced in V1:
* `TaskStatus` no longer holds the live fields of a task; it now serves as a read-only snapshot of the task's current state.
* Getter and setter methods have been introduced for accessing and modifying task fields.
* Task fields are no longer exported — please use the provided getters and setters instead.
* The `SetStatus` method no longer triggers `logger.OnTaskStatusUpdate`; if logging is needed, please call it explicitly.

---

* license: BSD3
* owners:
  * Rayark Inc. website: https://www.rayark.com
    * name: Yao Wei | github: @medicalwei | email: mwei@rayark.com
    * name: autumn31 | github: @autumn31 | email: ytli@rayark.com
    * name: Himawari Tachibana | github: @fieliapm | email: fieliapm@rayark.com
