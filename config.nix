{
  Logger = {
    Formatter = "json";
    Level = "info";
  };
  HTTP = {
    Addr = "127.0.0.1:8080";
  };
  Handlers = [
    {
      Method = "get";
      Path = "/api/v1/events";
      Type = "latests";
      Latests = {
        Format = "json";
        Inputs = [
          {
            Consumer = {
              Format = "json";
              Queue = {
                Type = "nsq";
                Nsq = {
                  Addr = "127.0.0.1:4150";
                  Topic = "tickers";
                  Channel = "platypus-tc-latests";
                  ConsumerBufferSize = 128;
                };
              };
            };
            Store = {
              Type = "memoryttl";
              MemoryTTL = {
                TTL = "24h";
				        Resolution = "1m";
			        };
			      };
            Key = ''{{.market}}|{{.symbolPair}}'';
          }
          {
            Consumer = {
              Format = "json";
              Queue = {
                Type = "nsq";
                Nsq = {
                  Addr = "127.0.0.1:4150";
                  Topic = "changes";
                  Channel = "platypus-tc-latests";
                  ConsumerBufferSize = 128;
                };
              };
            };
            Store = {
              Type = "memoryttl";
              MemoryTTL = {
                TTL = "24h";
				        Resolution = "1m";
			        };
            };
            Key = ''{{.period}}|{{.symbolPair}}'';
          }
        ];
        Wrap = ''{{"{"}}{{range $i, $e := $}}{{if $i}},{{end}}"{{- $e.Input.Consumer.Queue.Nsq.Topic -}}":{{- (printf "%s" $e.Events.JSON) -}}{{end}}{{"}"}}'';
      };
    }
    {
      Method = "get";
      Path = "/api/v1/events/stream";
      Type = "streams";
      Streams = {
        Format = "json";
        Inputs = [
          {
            Consumer = {
              Format = "json";
              Queue = {
                Type = "nsq";
                Nsq = {
                  Addr = "127.0.0.1:4150";
                  Topic = "tickers";
                  Channel = "platypus-tc-streams";
                  ConsumerBufferSize = 128;
                };
              };
            };
          }
          {
            Consumer = {
              Format = "json";
              Queue = {
                Type = "nsq";
                Nsq = {
                  Addr = "127.0.0.1:4150";
                  Topic = "changes";
                  Channel = "platypus-tc-streams";
                  ConsumerBufferSize = 128;
                };
              };
            };
          }
        ];
        Wrap = ''{"type":"{{- .Input.Consumer.Queue.Nsq.Topic -}}","payload":{{- (printf "%s" .Event.JSON) -}}}'';
        Writer = {
          ScheduleTimeout = "10s";
          Pool = {
            Workers = 128;
            QueueSize = 1024;
          };
        };
      };
    }
  ];
}
