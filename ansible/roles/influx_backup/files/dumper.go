package main

import (
  "encoding/json"
  "io/ioutil"
  "time"
  "flag"
  "fmt"
  "log"
  "os"
  "os/exec"
  "strings"
)

type DumpConfig struct {
  Dir string
  Aws string
  Host string
  Username string
  Password string
  Database string
  Series string
}

func main() {
  var path string
  flag.StringVar(&path, "config", "~/dumper.json", "Dump config (JSON)")
  flag.Parse()
  config := Config(path)
  dirname := DirName(config)
  DumpSeries(config, dirname)
  
  if config.Aws != "" {
    UploadToAws(config, dirname)
  }
}


func DumpSeries(self *DumpConfig, dirpath string) {
  t := time.Now()
  today := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())  
  yts := t.Unix() - (60*60*24)
  ytime := time.Unix(yts, 0)

  yesterday := fmt.Sprintf("%d-%02d-%02d", ytime.Year(), ytime.Month(), ytime.Day())  
  dir_path := fmt.Sprintf("%s/%s", self.Dir, dirpath)

  err := os.MkdirAll(dir_path, 0755)

  dump_path := fmt.Sprintf("%s/%s", dir_path, "dump.json")

  if err != nil {
    log.Fatal(err)
  }


 args_string := fmt.Sprintf("--database=%s --host=%s -out=%s --username=%s --password=%s --series=%s --start='%s' --end='%s'",
    self.Database,
    self.Host,
    dump_path,
    self.Username,
    self.Password,
    self.Series,
    yesterday,
    today)

  args := strings.Split(args_string, " ")

  output, err := exec.Command("influxdb-dump", args...).CombinedOutput()

  log.Printf("Output: %s", output)

  if err != nil {
    log.Fatal(err)
  }
}

func UploadToAws(self *DumpConfig, dirpath string) {
  dir_path := fmt.Sprintf("%s/%s", self.Dir, dirpath)
  aws_path := fmt.Sprintf("s3://%s/%s", self.Aws, dirpath)
  args := []string{"s3","sync", dir_path, aws_path}
  output, err := exec.Command("aws", args...).CombinedOutput()

  log.Printf("Output: %s", output)

  if err != nil {
    log.Fatal(err)
  }
}

func DirName(self *DumpConfig) string {
  var path string
  t := time.Now()
  path = fmt.Sprintf("%d/%02d/%02d", t.Year(), t.Month(), t.Day())
  return path
}


func Config(path string) *DumpConfig {
  b, err := ioutil.ReadFile(path)
  if err != nil {
    log.Fatalf("can't read config in file %q: %v", path, err)
  }
  var cfg DumpConfig
  if err := json.Unmarshal(b, &cfg); err != nil {
    log.Fatalf("error parsing config in file %q: %v", path, err)
  }
  return &cfg
}
