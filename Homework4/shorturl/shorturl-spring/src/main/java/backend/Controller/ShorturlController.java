package backend.Controller;

import backend.Service.ShorturlService;
import org.springframework.web.bind.annotation.*;


@RequestMapping
@RestController
public class ShorturlController {
    private final ShorturlService shorturlService;

    public ShorturlController(ShorturlService shorturlService) {
        this.shorturlService = shorturlService;
    }


    @RequestMapping(path = "/new")
    @ResponseBody
    public String createShort(@RequestParam String url){
        return shorturlService.createShorturl(url);
    }

    @RequestMapping(path = "/short/{data}")
    @ResponseBody
    public String getShort(@PathVariable String data){
        return shorturlService.getOriginal(data);
    }
}
