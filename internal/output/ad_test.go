package output

import (
	json "encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

var jsonstr = `{"status":1,"msg":"success","data":{"session_id":"5ae043c68423fd1b25a72154","parent_session_id":"5b3cbe40c4a18a61d807d0e0","ad_type":94,"template":2,"unit_size":"414x736","ads":[{"id":192797961,"title":"填色一二三：著色本 (Color by Number)","desc":"按著數字上色的新遊戲，肯定叫你愛不釋手！提供各式超有趣圖案任你選擇，只要依循數字上色，就能填出活靈活現的畫作！上色輕鬆無比，前所未見！\n\n色遊戲！\n\n......","package_name":"id1317978215","icon_url":"https://cdn-adn-https.rayjump.com/cdn-adn/dmp/18/01/04/12/28/5a4dad67b3918.jpg","image_url":"https://cdn-adn-https.rayjump.com/cdn-adn/v2/offersync/18/05/29/06/07/5b0c7db9b30bc.jpg","image_size":"VIDEO","impression_url":"https://sg01.rayjump.com/impression?k=5b3cbe40c4a18a61d87d0e37\u0026p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026x=0\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026al=\u0026csp=i%2BMeGUjPiUhPfA3FinEe6acIfnDBi%2BMeGUEA6aiM6aj%3D","video_url":"LdxThdi1WBKUH79wDkx/WktTJdSAWgzt4ku2Y+v/DFKwWFf3Y02tH79XJr9XinlXiahXiaiXin3XiA3XHrfrHUDeDn3BH7NbDAReDkR2G7ieiaJUfn53inReiai/Y+vT","video_length":30,"video_size":1982241,"video_resolution":"480x854","video_end_type":2,"playable_ads_without_video":1,"watch_mile":30,"ctype":1,"adv_imp":[],"ad_url_list":[],"t_imp":0,"adv_id":903,"click_url":"https://tracking.lenzmx.com/click?mb_pl=ios\u0026mb_nt=cb12537\u0026mb_campid=nx_fgff_cbn_hk_ios_vid_non_24393\u0026mb_creative_id=729453\u0026mb_ad_type=8\u0026mb_package=id995122577\u0026aff_sub=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026mb_subid=31225\u0026mb_idfa=9C2AA127-EDF0-4D01-895C-A841C2659D8B","notice_url":"https://sg01.rayjump.com/click?k=5b3cbe40c4a18a61d87d0e37\u0026p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026al=\u0026csp=i%2BMeGUjPiUhPfA3FinEe6acIfnDBi%2BMeGUEA6aiM6aj%3D\u0026notice=1","fca":2,"fcb":2,"template":9,"ad_source_id":1,"app_size":"29","click_mode":6,"rating":4.5,"landing_type":0,"ctatext":"安装","c_ct":60,"link_type":1,"guidelines":"","reward_amount":0,"reward_name":"","offer_type":0,"retarget_offer":2,"ttc":false,"ad_tracking":{"mute":["https://sg01.rayjump.com/trackv2?p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026type=reward_video\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026key=mute"],"unmute":["https://sg01.rayjump.com/trackv2?p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026type=reward_video\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026key=unmute"],"endcard_show":["https://sg01.rayjump.com/trackv2?p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026type=reward_video\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026key=endcard_show"],"close":["https://sg01.rayjump.com/trackv2?p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026type=reward_video\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026key=close"],"play_percentage":[{"rate":0,"url":"https://sg01.rayjump.com/trackv2?p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026type=reward_video\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026key=play_percentage\u0026rate=0"},{"rate":25,"url":"https://sg01.rayjump.com/trackv2?p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026type=reward_video\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026key=play_percentage\u0026rate=25"},{"rate":50,"url":"https://sg01.rayjump.com/trackv2?p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026type=reward_video\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026key=play_percentage\u0026rate=50"},{"rate":75,"url":"https://sg01.rayjump.com/trackv2?p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026type=reward_video\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026key=play_percentage\u0026rate=75"},{"rate":100,"url":"https://sg01.rayjump.com/trackv2?p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026type=reward_video\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026key=play_percentage\u0026rate=100"}],"pause":["https://sg01.rayjump.com/trackv2?p=fHx8fHx8fHJld2FyZGVkX3ZpZGVvfFZJREVPfHxpb3N8OS4zLjV8bWlfMi44LjF8aXBob25lNywxfDQxNHg3MzZ8MXx8emgtSGFucy1DTnx3aWZpfDQ2MDAyfHx8TU5vcm1hbEFscGhhTW9kZWxSYW5rZXI7MTA7MjE3OzIxOzA7MjE7MjsxX3VzZXJfaGlzdG9yeV9kbXBfdGFnXzI0X2hvdXItMl9jcGVfYmFzZS0zX2NyZWF0aXZlX3RyeV9uZXctNV9jY3BfYmFzZTswOzB8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHx8fHx8fHwzOS4xMDkuMTI0LjkzfHx8fHx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8fHx8MXx8MTg3MkM4QjMtNEMyRi00Mjk1LTlEREMtQjI3MjhEREIyMkIwLHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8fHx8\u0026q=a_i09M6dfgiaj%2FhrcPLg5whoPUYF2IfkRADFzQfaSUf7jeG7jFikN9fFNMHnib6a50iFf0HnNMDAxtinttfUc3GavbHaS3HgM9iAlT6aieiUR26aN2Gn5IGnvA6ajPiUhPfA3Fi%2BMeWURM6aj%2FiUSIDkx%2FH%2Bx6DkxAH%2BzFH%2BzIfbeRZbMwi%2BMF6acIiAjBiU5IiA3M6aR%2FfAl%2FideXh75%2FD%2BSu6aQI6o2IiASIidMAinRBf%2BMM6acIidMe6dMT6dMM67Q3Gn32inRBfnhb6acIiUN9inlFiUDTibMB6dMe6aSI6aSIideIideYRUj%2FiUv0WoReWURMR0M0iZRsRUj0%2B%2BM2faSIidMM6v%3D%3D\u0026type=reward_video\u0026r=eyJnaWQiOiI2MWZjN2I2Y2ZlZjgwZmY1MGJkNzcyNTI2NWFhMmUzNCIsInRwaWQiOjAsImNyYXQiOjgsImFkdl9jcmlkIjo3Mjk0NTMsImljYyI6MSwiZ2xpc3QiOiIxMDYsMjQxMDU4MjM5Miw3NDI3MTMsfDIwMSwyNDgxODYyNjQzLDcyOTQ1Myw0ODB4ODU0fDQwMSwyNDEwNTgyMzY5LDAsfDQwMiwyNDEwNTgyMzcwLDAsfDQwMywyNDEwNTgyMzcxLDAsfDQwNCwyNDEwNTgyMzcyLDAsfDQwNSwyNDEwNTgyMzczLDAsfDQwNiwyNDEwNTgyMzc0LDAsIiwicGkiOjEuMiwicG8iOjEuMiwiZGNvIjowfQ%3D%3D\u0026key=pause"]},"storekit":0,"md5_file":"","number_rating":33171,"icon_mime":"image/jpeg","image_mime":"image/jpeg","image_resolution":"1200x627","video_width":480,"video_height":854,"bitrate":527,"video_mime":"video/mp4","sub_category_name":["6016","7012","7014"],"storekit_time":2,"endcard_click_result":1,"c_toi":2,"imp_ua":1,"c_ua":1}],"html_url":"","end_screen_url":"https://hybird.rayjump.com/offerwall/tpl/mintegral/endscreen.v4.html?unit_id=4595\u0026sdk_version=mi_2.8.1","only_impression_url":"https://sg01.rayjump.com/onlyImpression?k=5b3cbe40c4a18a61d807d0df\u0026p=ODM4NHwzMTIyNXw0NTk1fDB8MHx8b3BlbmFwaXxyZXdhcmRlZF92aWRlb3wxMjh4MTI4fDd8aW9zfDkuMy41fG1pXzIuOC4xfGlwaG9uZTcsMXw0MTR4NzM2fDF8SEt8emgtSGFucy1DTnx3aWZpfDQ2MDAyfGFkbmV0X2Fkc2VydmVyfHxNTm9ybWFsQWxwaGFNb2RlbFJhbmtlcjsxMDsyMTc7MjE7MDsyMTsyOzFfdXNlcl9oaXN0b3J5X2RtcF90YWdfMjRfaG91ci0yX2NwZV9iYXNlLTNfY3JlYXRpdmVfdHJ5X25ldy01X2NjcF9iYXNlOzA7MHw1YjNjYmU0MGM0YTE4YTYxZDgwN2QwZGZ8NWIzY2JlNDBjNGExOGE2MWQ4MDdkMGRmfHwzMHwwfE1vemlsbGElMkY1LjAlMjAlMjhpUGhvbmUlM0IlMjBDUFUlMjBpUGhvbmUlMjBPUyUyMDlfM181JTIwbGlrZSUyME1hYyUyME9TJTIwWCUyOSUyMEFwcGxlV2ViS2l0JTJGNjAxLjEuNDYlMjAlMjhLSFRNTCUyQyUyMGxpa2UlMjBHZWNrbyUyOSUyME1vYmlsZSUyRjEzRzM2fHw1YjNjYmU0MGM0YTE4YTYxZDgwN2QwZGZ8MzkuMTA5LjEyNC45M3x8fHwxMy4yNTAuMjA5LjEzNXx8fHx8OUMyQUExMjctRURGMC00RDAxLTg5NUMtQTg0MUMyNjU5RDhCfDIuNzguMHxhcHBsZXx8NWFlMDQzYzY4NDIzZmQxYjI1YTcyMTU0fDViM2NiZTQwYzRhMThhNjFkODA3ZDBlMHx8MzkwfDF8MzEyMjV8MXwyfDE4NzJDOEIzLTRDMkYtNDI5NS05RERDLUIyNzI4RERCMjJCMCx8fHxbWyIxOTI3OTc5NjEiLCI5MDMiLCIiLCIiLCIxIiwiMSIsIjViM2NiZTQwYzRhMThhNjFkODdkMGUzNyJdXXx8MHx8fHxjb20ua2lsb28uc3Vid2F5c3VyZi5jbnx8MXx8NHwwfHxpZDk5NTEyMjU3N3wwfHx8fDB8MHx8fHx8fHwwfHx8fHx8MHx8fHwwfHx8TW96aWxsYSUyRjUuMCUyMCUyOE1hY2ludG9zaCUzQiUyMEludGVsJTIwTWFjJTIwT1MlMjBYJTIwMTBfMTNfMSUyOSUyMEFwcGxlV2ViS2l0JTJGNTM3LjM2JTIwJTI4S0hUTUwlMkMlMjBsaWtlJTIwR2Vja28lMjklMjBDaHJvbWUlMkY2Ni4wLjMzNTkuMTgxJTIwU2FmYXJpJTJGNTM3LjM2fHx8fHx8\u0026csp=i%2BMM6acIfnDBi%2BMeGUv1idMAidMe"}}`

var mvResult = &MobvistaResult{}

func BenchmarkMobvistaResultJsonMarshal(b *testing.B) {
	err := json.Unmarshal([]byte(jsonstr), mvResult)
	if err != nil {
		panic(err)
	}

	b.Run("encodingjson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = json.Marshal(mvResult)
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			_, _ = json.Marshal(mvResult)
		}
	})
}
