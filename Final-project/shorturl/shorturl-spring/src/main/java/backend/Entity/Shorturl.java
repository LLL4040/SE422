package backend.Entity;

import javax.persistence.*;

@Entity
@Table(name = "shorturl")
public class Shorturl {
    @Id
    @Column(name = "id")
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    @Column(name = "original")
    private String original;
    @Column(name = "shorturl")
    private String shorturl;

    public Shorturl() {
    }

    public Shorturl(String original, String shorturl) {
        this.original = original;
        this.shorturl = shorturl;
    }

    public void setOriginal(String original) {
        this.original = original;
    }

    public void setShorturl(String shorturl) {
        this.shorturl = shorturl;
    }

    public String getOriginal() {
        return original;
    }

    public String getShorturl() {
        return shorturl;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Long getId() {
        return id;
    }
}
