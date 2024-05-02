# Priority Scheduling in CPU Core Constrained Environment

This Go application has a priority scheduler with a time-based priority increment feature. The time-based priority-increment feature prevents starvation of tasks when there is an overwhelming number of high-priority tasks.

This application can take any level of priority. Send requests to `http:localhost:8000/<priority-level>`. The `priority-level` is a parameter is a number, the higher the value, the higher the priority.