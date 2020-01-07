package backend.Repository;


import backend.Entity.Shorturl;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;


@Repository
public interface ShorturlRepository extends JpaRepository<Shorturl, String> {
    Boolean existsByShorturl(String shorturl);

    Boolean existsByOriginal(String original);

    Shorturl findByOriginal(String original);

    Shorturl findByShorturl(String shorturl);
}
