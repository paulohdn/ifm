/*
SELECT fra_id, CONVERT(fra_nm using latin1), fra_ocor, CONVERT(fra_sig_pt USING latin1)
FROM frase
WHERE fra_nm REGEXP '^[^A-Za-z0-9]'
 --and NOT (ASCII(SUBSTRING(fra_nm, 1)) in (' ', '-', '_', '.', '"'))
ORDER BY fra_ocor DESC;
*/
select * from frase order by fra_ocor;
/*
SELECT fra_id, left(CONVERT(fra_nm using latin1),20), fra_ocor, left(CONVERT(fra_sig_pt USING latin1),20), concat('[',SUBSTRING(fra_nm, 1,1),']') car
FROM frase
WHERE (SUBSTRING(fra_nm, 1,1) in (' ', '-', '_', '.', ''))
ORDER BY fra_ocor DESC;
*/
/*
"./subtitles4/friends.s04e01.720p.bluray.x264-psychd.srt"
SELECT pal_id, convert(pal_nm using latin1), pal_ocor, CONVERT(pal_sig_pt USING latin1)
    FROM palavra
WHERE pal_nm REGEXP '^[^A-Za-z0-9]'
ORDER BY pal_ocor DESC;
GO
*/
/*
SELECT fra_id, fra_nm, fra_ocor, CONVERT(fra_sig_pt USING latin1)
FROM frase
WHERE NOT (ASCII(SUBSTRING(fra_nm, -1)) BETWEEN 48 AND 57) -- Check if not a digit (0-9)
   AND NOT (ASCII(SUBSTRING(fra_nm, -1)) BETWEEN 65 AND 90) -- Check if not an uppercase letter (A-Z)
   AND NOT (ASCII(SUBSTRING(fra_nm, -1)) BETWEEN 97 AND 122) -- Check if not a lowercase letter (a-z)
ORDER BY fra_ocor DESC;
*/
/*
SELECT pal_id, convert(pal_nm using latin1), pal_ocor, CONVERT(pal_sig_pt USING latin1)
    FROM palavra
WHERE NOT (ASCII(SUBSTRING(pal_nm, -1)) BETWEEN 48 AND 57) -- Check if not a digit (0-9)
   AND NOT (ASCII(SUBSTRING(pal_nm, -1)) BETWEEN 65 AND 90) -- Check if not an uppercase letter (A-Z)
   AND NOT (ASCII(SUBSTRING(pal_nm, -1)) BETWEEN 97 AND 122) -- Check if not a lowercase letter (a-z)
ORDER BY pal_ocor DESC;
*/


/*
SELECT df.dob_id, df.leg_id, f.fra_nm from frase f 
inner join dob_fra df on df.fra_id=f.fra_id 
where f.fra_nm like '%h6y%'
GO
*/
