// Package goblast provides a highly interface for accessing the NCBI BLAST API.
//
// API documentation can be found here: https://ncbi.github.io/blast-cloud/dev/api.html
//
// Code Example:
//
// TODO: example code
/*
QUERY	Search query	String	Put *	Accession, GI, or FASTA.
DATABASE	BLAST database	String	Put *	Database from appendix 2 or one uploaded to blastdb_custom (see appendix 4)
PROGRAM	BLAST program	String	Put *	One of blastn, megablast, blastp, blastx, tblastn, tblastx
FILTER	Low complexity filtering	String	Put	F to disable. T or L to enable. Prepend “m” for mask at lookup (e.g., mL)
FORMAT_TYPE	Report type	String	Put, Get	HTML, Text, XML, XML2, JSON2, or Tabular. HTML is the default.
EXPECT	Expect value	Double	Put	Number greater than zero.
NUCL_REWARD	Reward for matching bases (BLASTN and megaBLAST)	Integer	Put	Integer greater than zero.
NUCL_PENALTY	Cost for mismatched bases (BLASTN and megaBLAST)	Integer	Put	Integer less than zero.
GAPCOSTS	Gap existence and extension costs	String	Put	Pair of positive integers separated by a space such as “11 1”.
MATRIX	Scoring matrix name	String	Put	One of BLOSUM45, BLOSUM50, BLOSUM62, BLOSUM80, BLOSUM90, PAM250, PAM30 or PAM70. Default: BLOSUM62 for all applicable programs
HITLIST_SIZE	Number of databases sequences to keep	Integer	Put,Get	Integer greater than zero.
DESCRIPTIONS	Number of descriptions to print (applies to HTML and Text)	Integer	Put,Get	Integer greater than zero.
ALIGNMENTS	Number of alignments to print (applies to HTML and Text)	Integer	Put,Get	Integer greater than zero.
NCBI_GI	Show NCBI GIs in report	String	Put, Get	T or F
RID	BLAST search request identifier	String	Get *, Delete *	The Request ID (RID) returned when the search was submitted
THRESHOLD	Neighboring score for initial words	Integer	Put	Positive integer (BLASTP default is 11). Does not apply to BLASTN or MegaBLAST).
WORD_SIZE	Size of word for initial matches	Integer	Put	Positive integer.
COMPOSITION_BASED_STATISTICS	Composition based statistics algorithm to use	Integer	Put	One of 0, 1, 2, or 3. See comp_based_stats command line option in the BLAST+ user manual for details.
FORMAT_OBJECT	Object type	String	Get	SearchInfo (status check) or Alignment (report formatting).
*/
package goblast
