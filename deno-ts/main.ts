import { BufReader } from "https://deno.land/std@0.174.0/io/buffer.ts";

function main() {
  const listener = Deno.listen({
    port: 8080,
    hostname: '0.0.0.0',
    transport: 'tcp',
  })

  console.log(`Listening on port 8080`);

  httpHandler(listener);
}

async function httpHandler(listener: Deno.Listener) {
  for await (const conn of listener) {
    const bufReader = new BufReader(conn);
    const lines = await bufReader.readLine();
    if (!lines) break;

    const [method, path, protocol, ...rest] = new TextDecoder().decode(lines.line).split(' ');

    console.log(`Akses ${path} dengan method ${method} menggunakan protocol ${protocol}`)

    if (method === 'GET') {
      await Promise.all([
        conn.write(new TextEncoder().encode("HTTP/1.1 200 OK\r\n")),
        conn.write(new TextEncoder().encode("Content-Type: application/json\r\n")),
        conn.write(new TextEncoder().encode("\r\n")),
        conn.write(new TextEncoder().encode(`
          {
            "message": "Anda mengakses ${path} dengan method ${method}",
            "data": "Disini data"
          }
        `))
      ])
    }

    conn.close();
  }
}

// run application
main();

