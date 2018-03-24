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
            Format = "json";
            Consumer = {
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
            Key = ''{{.Message.market}}|{{.Message.symbolPair}}'';
          }
          {
            Format = "json";
            Consumer = {
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
            Key = ''{{.Message.period}}|{{.Message.symbolPair}}'';
          }
        ];
        Wrap = ''{{"{"}}{{range $i, $e := $}}{{if $i}},{{end}}"{{- $e.Config.Consumer.Queue.Nsq.Topic -}}":{{- (printf "%s" $e.Data) -}}{{end}}{{"}"}}'';
      };
    }
    {
      Method = "get";
      Path = "/api/v1/events/stream";
      Type = "streams";
      Streams = {
        Inputs = [
          {
            Consumer = {
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
        Wrap = ''{"type":"{{- .Config.Consumer.Queue.Nsq.Topic -}}","payload":{{- (printf "%s" .Message) -}}}'';
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
