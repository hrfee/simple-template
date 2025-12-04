import { Template } from "./dist/index.js";
let tmpl = process.argv[2];
let mapBase = JSON.parse(process.argv[3]);
let b = new Map();
if (mapBase) {
    for (let key of Object.keys(mapBase)) {
        b.set(key, mapBase[key]);
    }
}
let [out, err] = Template(tmpl, b)
process.stdout.write(out);
if (err == null) out = `nil`;
else out = err.name;
process.stderr.write(out);
