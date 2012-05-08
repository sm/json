/*

  NAME

  sm-json - find variable in json file

  SYNOPSIS

  Usage of ./sm-mustache:
    -path="{path}": path to value
    -uri="{uri}": json uri or path to file

    -uri|--uri "{uri}"
      
      specify the json file to use

    -path|--path "{path}"

      specify the path to variable
   
  DESCRIPTION

  sm-json will read the given json file and get value. 

*/


package main

import (
  "encoding/json"
  "strings"
  "flag"
  "log"
  "io/ioutil"
  "bufio"
  "os"
  "fmt"
)

var path  = flag.String("path","{path}", "path of variable")
var uri   = flag.String("uri","{uri}", "json uri or path to file")
var f map[string]interface{}
var value string

func find_value(key []string, data map[string]interface{}, index int) string {
  switch v := data[key[index]].(type) {
  case string:
    return v
  case map[string]interface{}:
    return find_value(key,v,index+1)
  default:
    log.Fatal("ERROR: Error with value ", v)
  }
  return ""
}

func find_map_value(data map[string]interface{}) string {
  var my_val string
  for key, val := range data {
    my_val += key + " " + val.(string) + " "
  }
  return my_val
}

func main() {

  flag.Parse()

  if (*uri == "{uri}") {
    log.Fatalf("ERROR: A template file must be given! %s",*uri)
  } else {
    file, err := os.Open(*uri)
    if err != nil {
      log.Fatal("ERROR FILE NOT REAL")
    }
    bfile := bufio.NewReader(file)
    buf, err := ioutil.ReadAll(bfile)
    if file == nil && err != nil {
      log.Fatalf("ERROR: The json file %s does not exist!",*uri)
    }

    json_err := json.Unmarshal(buf,&f)

    if json_err != nil {
      log.Fatal("ERROR: unable to read json file")
    }

    //Parse path
    path_array := strings.Split(*path,"/")
    for i := range path_array {
      switch v := f[path_array[i]].(type) {
      case string:
        if(len(path_array) == i+1) {
          value = v
        } else {
          log.Fatal("ERROR: Found value too soon")
        }
      case map[string]interface{}:
        if(len(path_array) == i+1) {
          value = find_map_value(v)
        } else {
          value = find_value(path_array,v,i+1)
        }
        break
      case []interface{}:
        // @TODO make implementation better
        for index, data := range v {
          switch mapped_data := data.(type) {
          case map[string]interface{}:
              value += find_map_value(mapped_data)
          default:
            log.Fatal("ERROR: Data is not correct on index ", index)
          }
        }
      default:
      }
      if(len(value) > 0) {
        break
      }
    }

    fmt.Println(value)
  }
}
