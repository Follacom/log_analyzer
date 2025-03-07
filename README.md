# Log Analyzer

*Parse your logs to keep it clean*

## Apache HTTP Setup (>2.4.58 only)

There is a specific configuration to add to your Apache HTTP configuration file.

Log format should be structured like so:
```
LogFormat "[a:\"%a\"] [A:\"%A\"]" custom_access
ErrorLogFormat "[a:\"%a\"] [A:\"%A\"] [E:\"%E\"] [F:\"%F\"] [k:\"%k\"] [l:\"%l\"] [L:\"%L\"] [cL:\"%{c}L\"] [m:\"%m\"] [M:\"%M\"] [P:\"%P\"] [T:\"%T\"] [gT:\"%{g}T\"] [t:\"%{%Y-%m-%d %T}t.%{usec_frac}t %{%z}t\"] [v:\"%v\"]"
```