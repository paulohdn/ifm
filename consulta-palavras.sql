/*
SELECT 	DISTINCT 
    p.pal_id,
	p.pal_nm,
	CONVERT(p.pal_sig_pt USING latin1) pal_sig_pt,
	p.pal_ocor
FROM dob_pal dp 
INNER JOIN palavra p ON dp.pal_id = p.pal_id 
--where dp.dob_id = 234 and p.pal_id not in (select pal_id from dob_pal where dob_id < 234)
	ORDER BY p.pal_ocor DESC 
;
*/

SELECT fra_id, fra_nm, fra_ocor, CONVERT(fra_sig_pt using latin1) fra_sig_pt
    FROM frase
GO


--- CONVERT(field_name USING latin1)
--- collate utf8mb4_bin