package backend.Service;

import backend.Entity.Shorturl;
import backend.Repository.ShorturlRepository;
import org.springframework.stereotype.Service;


@Service
public class ShorturlService {
    private final ShorturlRepository shorturlRepository;
    private static final char[] digits = {'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
            'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
            'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
            'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
            'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'};

    public ShorturlService(ShorturlRepository shorturlRepository) {
        this.shorturlRepository = shorturlRepository;
    }

    public String createShorturl(String original) {
        Shorturl shorturl = new Shorturl();
        shorturl.setOriginal(original);
        shorturl = shorturlRepository.save(shorturl);
        long id = shorturl.getId();
        String shorturlString = toOtherNumberSystem(id, 62);
        shorturl.setShorturl(shorturlString);
        shorturlRepository.save(shorturl);
        return shorturlString;
    }

    public String getOriginal(String shorturl) {
        return shorturlRepository.findByShorturl(shorturl).getOriginal();
    }

    public static String toOtherNumberSystem(long number, int seed) {
        if (number < 0) {
            number = ((long) 2 * 0x7fffffff) + number + 2;
        }
        char[] buf = new char[32];
        int charPos = 32;
        while ((number / seed) > 0) {
            buf[--charPos] = digits[(int) (number % seed)];
            number /= seed;
        }
        buf[--charPos] = digits[(int) (number % seed)];
        return new String(buf, charPos, (32 - charPos));
    }
}
