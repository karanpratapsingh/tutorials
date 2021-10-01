import cluster, { Worker } from 'cluster';
import os from 'os';
import * as Server from './server';

if (cluster.isMaster) {
  const cores = os.cpus().length;

  console.log(`Total cores: ${cores}`);
  console.log(`Primary process ${process.pid} is running`);

  for (let i = 0; i < cores; i++) {
    cluster.fork();
  }

  cluster.on('exit', (worker: Worker, code) => {
    console.log(`Worker ${worker.process.pid} exited with code ${code}`);
    console.log('Fork new worker!');
    cluster.fork();
  });
} else {
  Server.start();
}
