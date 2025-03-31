package constants

const (
	Parser_prompt = `context: you are intelligent BDD tester having knowledge of karate testing framework
                     given the following input json data: %s in which the background field has karate feature
					 file background related data , and the scenario fields are list of scenarios that
                     every one of which will map to one karate scenario.And the attached filenames are given in %s list.
                     
					 your task is to read the above mentioned input json data understand it and intelligently parse the 
					 data in accordance to this output json file :%s. please fill this output json with the extracted
                     values and map to the fields of output json file intelligently.If any additional field are there
                     in the input json , you can add those details in additional_details instructions fields as an list
                     of strings, the request file location and header file location can be dynamically set if their respective
					 saving path is given so we can intelligently chose the name from the filenames given and the requestpath should be.
                     request_file_save_path+'/'+request_file_name.json. If the user has provided the header file and the path for it 
 					 to save then header_file_save_path +'/'+header_file_name.json else it should be null. dont assume any path by yourself

                     DISCLAIMER: Alert !! Stricly only send the output in json format and follow the output json format
					 never add any extra fields in that json, and dont assume anything on your own, No explanation , warning
					 nothing extra noise needed only send the updated json format 
					`
	Gemini_url     = `https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s`
	Gemini_Request = `{
                       "contents": {
                         "role": "user",
                         "parts": [
                           {
                             "text": "%s"
                           }
                         ]
                        }
                  }`
)
