package model_convert

func GenerateHTTP(replacement ... map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	var functionName string
	var payload string
	var _payload string
	_, _, _ = functionName, payload, _payload

	// TODO
	rs := `
// Auto-generate by github.com/fwhezfwhez/model_convert.GenerateHTTP
// TODO:
// - /path/to/${http_client}
func ${function_name}(${payload}) error {
    url := "${url}"
    method := "${method}"
    // TODO: add payload
    payload := map[string]interface{} {
        ${_payload}
    }
    req, e:= http.NewRequest(method, url, bytes.NewReader(buf))
    if e!=nil {
        return errorx.Wrap(e)
    }
    
    resp ,e:=${http_client}.Do(req)
    if e!=nil {
        return errorx.Wrap(e)
    }

    if resp == nil {
        return errorx.NewFromString("resp nil")
    }
    if resp.Body == nil {
        return errorx.NewFromString("resp.body nil")
    }

    if resp.StatusCode != 200 {
        b,e := ioutil.ReadAll(resp.Body)
        if e!=nil {
            return errorx.Wrap(e)
        }
        return errorx.NewFromStringf("recv status '%d', body '%s'", resp.StatusCode, string(b))
    }    
    
    ${response_struct}

    var ${_response_instance} ${_response_struct_name}
    b, e:= ioutil.ReadAll(resp.Body)
    if e!=nil {
        return errorx.Wrap(e)
    }
    e = json.Unmarshal(b, &${_response_instance})
    if e!=nil {
        return errorx.Wrap(e)
    }
    
    // TODO, define fail and success to ${_response_instance}
    return nil
}
`
	return rs
}
