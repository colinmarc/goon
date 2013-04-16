It's a baby language!

    if a:
      DoThing()
    elif b:
      DoOtherThing
    else:
      DoThirdThing()

    DoFinalThing() if b

    forever:
      DoThing()

    DoThing() forever

    for word in words:
      DoThing(word)

    DoThing(word) for word in words

methods, blocks... concurrency!

    ParallelSearch (terms) ->
      results = terms.map ->
        # launches in the background
        res = search(term)...
        results.append(res)

      # results is a list of promises, but you don't really care.
      return results

promises are kinda cool

    promise = search('foobar')...

    print ...promise # this'll block until the promise is fulfilled

    #maybe lists know how to wait on all their promises?
    results = ParallelSearch(['foo', 'bar'])
    return ...results

all blocks are just generators that restart when you call them again. they
can also return multiple times

    Counter ->
      i = 0
      forever:
        return i++

    Total (inc) ->
      i = 0
      return i + inc forever=

    # the new keyword 'forks' an instance

    a = new Counter
    b = new Counter

    i = a() # 1
    j = b() # 1

    c = new b
    k = c() # 2

variables inside blocks are available via dot notation. you can use this to
create classes

    Car ->
      speed = 0

      accelerate ->
        speed++

      stopped ->
        return (speed == 0)

    c = new Car() # we're calling the 'static' Car, so everything is
                  # initialized, then forking

    c.speed # 0
    c.stopped() # true

    c.accelerate() # 1

    c.speed # 1
    c.stopped() # false

inheritance, maybe?

    SSLServer (host, port, sslopts) ->
      extend new Server(host, port)
      # etc

a server!

    Server (host, port) ->
      sock = new Socket(host, port)

      Listen (host, port) ->
        sock.Bind(host, port)
        sock.Listen()

      Accept ->
        return sock.accept()

    HTTPStream (chunk) ->
      req = ''
      finished = false

      forever:
        req += chunk
        if chunk.Empty() or IsFullHTTPRequest(req):
          finished = true
          return new HTTPRequest(req)
          req = ''
        else:
          return

    HandleConnection (sock) ->
      stream = new HTTPStream # no parens here - this is a generator and doesn't
                              # initialize until it's called for the first time
      response = 'Hello World!'

      forever:
        chunk = sock.Read(1024)
        if chunk.Empty():
          sock.WriteAll(response)
          break

        req = stream(chunk)
        if req:
          sock.WriteAll(response)
          break

    server = new Server('', 9599)
    forever:
      req = server.Accept()
      HandleConnection(req)...