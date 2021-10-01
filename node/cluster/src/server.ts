import express, { Request, Response } from 'express';

export function start(): void {
  const app = express();

  app.get('/api/intense', (req: Request, res: Response): void => {
    console.time('intense');
    intenseWork();
    console.timeEnd('intense');
    res.send('Done!');
  });

  app.listen(4000, () => {
    console.log(`Server started with worker ${process.pid}`);
  });
}

/**
 * Mimics some intense server-side work
 */
function intenseWork(): void {
  const list = new Array<number>(1e7);

  for (let i = 0; i < list.length; i++) {
    list[i] = i * 12;
  }
}
