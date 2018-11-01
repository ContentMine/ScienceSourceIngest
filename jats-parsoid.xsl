<?xml version="1.0"?>

<xsl:stylesheet version="1.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    xmlns:mml="http://www.w3.org/1998/Math/MathML"
    exclude-result-prefixes="xlink mml">

    <xsl:include href="jats-common.xsl"/>

    <xsl:output method="html" doctype-system="about:legacy-compat"/>

    <xsl:template match="/">
        <html>
            <xsl:call-template name="make-html-header"/>
            <body>
                <xsl:apply-templates/>
            </body>
        </html>
    </xsl:template>

    <xsl:template name="make-html-header">
        <head>
            <title>
                <xsl:text>Converted JATS paper:</xsl:text>
            </title>
        </head>
    </xsl:template>

    <xsl:template match="article">
        <xsl:call-template name="make-article"/>
    </xsl:template>

    <xsl:template name="make-article">
        <xsl:apply-templates select="front"/>
        <xsl:for-each select="body">
            <h2 id="Paper"><xsl:text>Paper</xsl:text></h2>
            <xsl:apply-templates/>
        </xsl:for-each>
        <xsl:apply-templates select="back"/>
    </xsl:template>

    <xsl:template match="front">
        <section>
            <h2 id="Journal_Information"><xsl:text>Journal Information</xsl:text></h2>
            <xsl:for-each select="journal-meta">
                <xsl:apply-templates/>
            </xsl:for-each> <!-- journal-meta -->
        </section>

        <section>
            <xsl:for-each select="article-meta">
                <xsl:apply-templates/>
            </xsl:for-each> <!-- article-meta -->
        </section>
    </xsl:template>

    <xsl:template match="back">
        <xsl:apply-templates select="ack | ref-list"/>
    </xsl:template>

</xsl:stylesheet>
