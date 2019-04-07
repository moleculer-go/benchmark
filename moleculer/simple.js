const { ServiceBroker } = require("moleculer");
const broker = new ServiceBroker();

broker.createService({
  name: "math",
  actions: {
    add({ params }) {
      return params.a + params.b;
    }
  }
});

function add(a, b) {
  return a + b;
}

async function benchmark(bench) {
  const { total, rps, avg } = await doForOneSecond(async () => {
    //const r = await broker.call("math.add", { a: 5, b: 3 });
    // return r === 8;
    const r = add(5, 3);
    return r === 8;
  });
  bench.save("rps", rps);
  bench.save("total", total);
  bench.save("avg", avg);
}

broker.start().then(() => runTimes(benchmark, 10));

//TODO Move to another file !

const NANOS_P_SEC = BigInt(1e9);

async function doForOneSecond(work) {
  let total = 0;
  const start = process.hrtime.bigint();
  while (true) {
    const validCycle = await work();
    if (!validCycle) {
      new Error("invalid result");
    }
    total++;

    if (total % 1000 == 0) {
      const nano = process.hrtime.bigint() - start;
      const duration = Number(nano / NANOS_P_SEC);
      if (duration >= 1) {
        return {
          total,
          rps: total / duration,
          avg: duration / total,
          duration
        };
      }
    }
  }
}

async function runTimes(fn, times) {
  const bench = {
    values: {},
    save: (name, value) => {
      if (!bench.values[name]) {
        bench.values[name] = new Array();
      }
      bench.values[name].push(value);
    }
  };
  for (let i = 0; i < times; i++) {
    await fn(bench);
  }
  printBench(bench.values);
}

function stats(list) {
  let min = -1,
    max = 0,
    avg = 0,
    total = 0;
  let i = 0;
  for (; i < list.length; i++) {
    const item = list[i];
    if (item < min || min === -1) {
      min = item;
    }
    if (item > max) {
      max = item;
    }
    total += item;
  }
  avg = total / i;
  return { min, max, avg, total };
}

function printBench(bench) {
  for (key in bench) {
    const list = bench[key];
    const { min, max, avg, total } = stats(list);
    console.log(`\n **** Stats for [ ${key} ] **** `);
    console.log("\n avg: ", avg);
    console.log("\n min: ", min);
    console.log("\n max: ", max);
    console.log("\n total: ", total);
  }
}
