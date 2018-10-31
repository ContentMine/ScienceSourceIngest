<?xml version="1.0"?>

<xsl:stylesheet version="1.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    xmlns:mml="http://www.w3.org/1998/Math/MathML"
    exclude-result-prefixes="xlink mml">

    <xsl:include href="jats-common.xsl"/>

    <xsl:output method="text"/>

    <xsl:template match="/">
        <xsl:apply-templates/>
    </xsl:template>

    <xsl:template match="article">
        <xsl:call-template name="make-article"/>
    </xsl:template>

    <xsl:template name="make-article">
        <xsl:apply-templates/>
    </xsl:template>

    <xsl:template match="front" select="article-meta">
        <xsl:apply-templates/>
    </xsl:template>

    <xsl:template match="back">
    </xsl:template>

</xsl:stylesheet>
